package controllers

import (
    "github.com/gofiber/fiber/v2"
    "ticketing-backend/config"
    "ticketing-backend/models"
    "github.com/google/uuid"
)

type CreateTicketRequest struct {
    EventID          string `json:"event_id"`
    TicketCategoryID string `json:"ticket_category_id"`
    Quantity         int    `json:"quantity"`
}

func CreateTicket(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)

    var req CreateTicketRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Check ticket category
    var ticketCategory models.TicketCategory
    if err := config.DB.Where("ticket_category_id = ?", req.TicketCategoryID).First(&ticketCategory).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Ticket category not found",
        })
    }

    // Check quota
    if ticketCategory.Sold+req.Quantity > ticketCategory.Quota {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Not enough tickets available",
        })
    }

    // Create tickets
    var tickets []models.Ticket
    for i := 0; i < req.Quantity; i++ {
        ticket := models.Ticket{
            EventID:          req.EventID,
            TicketCategoryID: req.TicketCategoryID,
            OwnerID:          userID,
            Code:             uuid.New().String(),
        }
        tickets = append(tickets, ticket)
    }

    if err := config.DB.Create(&tickets).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create tickets",
        })
    }

    // Update sold count
    config.DB.Model(&ticketCategory).Update("sold", ticketCategory.Sold+req.Quantity)

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Tickets created successfully",
        "tickets": tickets,
    })
}

func GetTickets(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)

    var tickets []models.Ticket
    if err := config.DB.Preload("Event").Preload("TicketCategory").Where("owner_id = ?", userID).Find(&tickets).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch tickets",
        })
    }

    return c.JSON(fiber.Map{
        "tickets": tickets,
    })
}

func GetTicket(c *fiber.Ctx) error {
    ticketID := c.Params("id")
    userID := c.Locals("userID").(string)

    var ticket models.Ticket
    if err := config.DB.Preload("Event").Preload("TicketCategory").Where("ticket_id = ? AND owner_id = ?", ticketID, userID).First(&ticket).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Ticket not found",
        })
    }

    return c.JSON(fiber.Map{
        "ticket": ticket,
    })
}

func CheckInTicket(c *fiber.Ctx) error {
    ticketID := c.Params("id")

    var ticket models.Ticket
    if err := config.DB.Where("ticket_id = ?", ticketID).First(&ticket).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Ticket not found",
        })
    }

    if ticket.Status == "used" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Ticket already used",
        })
    }

    if err := config.DB.Model(&ticket).Update("status", "used").Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to check in ticket",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Ticket checked in successfully",
    })
}