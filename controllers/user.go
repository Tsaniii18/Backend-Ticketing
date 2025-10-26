package controllers

import (
    "github.com/gofiber/fiber/v2"
    "ticketing-backend/config"
    "ticketing-backend/models"
)

func GetProfile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)

    var user models.User
    if err := config.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    return c.JSON(fiber.Map{
        "user": user,
    })
}

type UpdateProfileRequest struct {
    Name      string `json:"name"`
    ProfilePic string `json:"profile_pic"`
}

func UpdateProfile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)

    var req UpdateProfileRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if err := config.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(models.User{
        Name:       req.Name,
        ProfilePic: req.ProfilePic,
    }).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update profile",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Profile updated successfully",
    })
}

func GetUsers(c *fiber.Ctx) error {
    var users []models.User
    if err := config.DB.Find(&users).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch users",
        })
    }

    return c.JSON(fiber.Map{
        "users": users,
    })
}

func VerifyUser(c *fiber.Ctx) error {
    userID := c.Params("id")

    var user models.User
    if err := config.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    if user.Role != "eo" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Only EO accounts can be verified",
        })
    }

    if err := config.DB.Model(&models.User{}).Where("user_id = ?", userID).Update("register_status", "approved").Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to verify user",
        })
    }

    return c.JSON(fiber.Map{
        "message": "User verified successfully",
    })
}