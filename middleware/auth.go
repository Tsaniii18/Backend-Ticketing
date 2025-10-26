package middleware

import (
    "strings"
    "github.com/gofiber/fiber/v2"
    "ticketing-backend/utils"
)

func AuthMiddleware(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Authorization header required",
        })
    }

    tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
    claims, err := utils.ValidateJWT(tokenString)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }

    c.Locals("userID", claims.UserID)
    c.Locals("role", claims.Role)
    return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {
    role := c.Locals("role").(string)
    if role != "admin" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Admin access required",
        })
    }
    return c.Next()
}

func EOMiddleware(c *fiber.Ctx) error {
    role := c.Locals("role").(string)
    if role != "eo" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Event Organizer access required",
        })
    }
    return c.Next()
}