package controllers

import (
    "github.com/gofiber/fiber/v2"
    "ticketing-backend/config"
    "ticketing-backend/models"
    "time"
)

type CreateEventRequest struct {
    Name        string    `json:"name"`
    DateStart   time.Time `json:"date_start"`
    DateEnd     time.Time `json:"date_end"`
    Location    string    `json:"location"`
    Description string    `json:"description"`
    Image       *string   `json:"image"`
    Flyer       *string   `json:"flyer"`
    Category    string    `json:"category"`
}

func CreateEvent(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    role := c.Locals("role").(string)

    if role != "eo" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Only EO can create events",
        })
    }

    var req CreateEventRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    event := models.Event{
        OwnerID:     userID,
        Name:        req.Name,
        DateStart:   req.DateStart,
        DateEnd:     req.DateEnd,
        Location:    req.Location,
        Description: req.Description,
        Image:       req.Image,
        Flyer:       req.Flyer,
        Category:    req.Category,
        Status:      "pending",
    }

    if err := config.DB.Create(&event).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create event",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Event created successfully",
        "event":   event,
    })
}

func GetEvents(c *fiber.Ctx) error {
    var events []models.Event
    if err := config.DB.Preload("Owner").Where("status = ?", "approved").Find(&events).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch events",
        })
    }

    return c.JSON(fiber.Map{
        "events": events,
    })
}

func GetEvent(c *fiber.Ctx) error {
    eventID := c.Params("id")

    var event models.Event
    if err := config.DB.Preload("Owner").Where("event_id = ?", eventID).First(&event).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Event not found",
        })
    }

    return c.JSON(fiber.Map{
        "event": event,
    })
}

func UpdateEvent(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    eventID := c.Params("id")

    var event models.Event
    if err := config.DB.Where("event_id = ?", eventID).First(&event).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Event not found",
        })
    }

    if event.OwnerID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "You can only update your own events",
        })
    }

    var req CreateEventRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if err := config.DB.Model(&event).Updates(models.Event{
        Name:        req.Name,
        DateStart:   req.DateStart,
        DateEnd:     req.DateEnd,
        Location:    req.Location,
        Description: req.Description,
        Image:       req.Image,
        Flyer:       req.Flyer,
        Category:    req.Category,
    }).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update event",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Event updated successfully",
    })
}

func VerifyEvent(c *fiber.Ctx) error {
    eventID := c.Params("id")

    var event models.Event
    if err := config.DB.Where("event_id = ?", eventID).First(&event).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Event not found",
        })
    }

    if err := config.DB.Model(&event).Update("status", "approved").Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to verify event",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Event verified successfully",
    })
}

func DeleteEvent(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    eventID := c.Params("id")

    var event models.Event
    if err := config.DB.Where("event_id = ?", eventID).First(&event).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Event not found",
        })
    }

    if event.OwnerID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "You can only delete your own events",
        })
    }

    if err := config.DB.Delete(&event).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to delete event",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Event deleted successfully",
    })
}