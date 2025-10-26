package controllers

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    "ticketing-backend/config"
    "ticketing-backend/models"
    "github.com/google/uuid"
)

type AddToCartRequest struct {
    TicketCategoryID string `json:"ticket_category_id"`
    Quantity         int    `json:"quantity"`
}

func AddToCart(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)

    var req AddToCartRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Check if ticket category exists
    var ticketCategory models.TicketCategory
    if err := config.DB.Where("ticket_category_id = ?", req.TicketCategoryID).First(&ticketCategory).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Ticket category not found",
        })
    }

    // Check if item already in cart
    var existingCart models.Cart
    if err := config.DB.Where("user_id = ? AND ticket_category_id = ?", userID, req.TicketCategoryID).First(&existingCart).Error; err == nil {
        // Update quantity if exists
        existingCart.Quantity += req.Quantity
        if err := config.DB.Save(&existingCart).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to update cart",
            })
        }
        return c.JSON(fiber.Map{
            "message": "Cart updated successfully",
        })
    }

    cart := models.Cart{
        UserID:           userID,
        TicketCategoryID: req.TicketCategoryID,
        Quantity:         req.Quantity,
    }

    if err := config.DB.Create(&cart).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to add to cart",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Item added to cart successfully",
    })
}

func UpdateCart(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)

    var req AddToCartRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    var cart models.Cart
    if err := config.DB.Where("user_id = ? AND ticket_category_id = ?", userID, req.TicketCategoryID).First(&cart).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Item not found in cart",
        })
    }

    if req.Quantity <= 0 {
        // Remove item if quantity is 0 or negative
        if err := config.DB.Delete(&cart).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to remove item from cart",
            })
        }
        return c.JSON(fiber.Map{
            "message": "Item removed from cart",
        })
    }

    cart.Quantity = req.Quantity
    if err := config.DB.Save(&cart).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update cart",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Cart updated successfully",
    })
}

func DeleteFromCart(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    ticketCategoryID := c.Query("ticket_category_id")

    if ticketCategoryID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Ticket category ID is required",
        })
    }

    var cart models.Cart
    if err := config.DB.Where("user_id = ? AND ticket_category_id = ?", userID, ticketCategoryID).First(&cart).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Item not found in cart",
        })
    }

    if err := config.DB.Delete(&cart).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to remove item from cart",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Item removed from cart successfully",
    })
}

func Checkout(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)

    // Get user's cart items
    var cartItems []models.Cart
    if err := config.DB.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch cart items",
        })
    }

    if len(cartItems) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Cart is empty",
        })
    }

    // Process checkout in transaction
    err := config.DB.Transaction(func(tx *gorm.DB) error {
        for _, item := range cartItems {
            // Get ticket category details
            var ticketCategory models.TicketCategory
            if err := tx.Where("ticket_category_id = ?", item.TicketCategoryID).First(&ticketCategory).Error; err != nil {
                return err
            }

            // Check quota
            if ticketCategory.Sold + item.Quantity > ticketCategory.Quota {
                return fiber.NewError(fiber.StatusBadRequest, "Not enough tickets available for category: " + ticketCategory.Description)
            }

            // Create tickets
            for i := 0; i < item.Quantity; i++ {
                ticket := models.Ticket{
                    EventID:          ticketCategory.EventID,
                    TicketCategoryID: item.TicketCategoryID,
                    OwnerID:          userID,
                    Code:             uuid.New().String(),
                }
                if err := tx.Create(&ticket).Error; err != nil {
                    return err
                }
            }

            // Update sold count
            if err := tx.Model(&models.TicketCategory{}).Where("ticket_category_id = ?", item.TicketCategoryID).Update("sold", ticketCategory.Sold + item.Quantity).Error; err != nil {
                return err
            }
        }

        // Clear cart
        if err := tx.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Checkout failed: " + err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "message": "Checkout successful",
    })
}