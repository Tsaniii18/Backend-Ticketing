package controllers

import (
    "github.com/gofiber/fiber/v2"
    "golang.org/x/crypto/bcrypt"
    "ticketing-backend/config"
    "ticketing-backend/models"
    "ticketing-backend/utils"
)

type RegisterRequest struct {
    Username               string  `json:"username"`
    Name                   string  `json:"name"`
    Email                  string  `json:"email"`
    Password               string  `json:"password"`
    Role                   string  `json:"role"`
    Organization           *string `json:"organization,omitempty"`
    OrganizationType       *string `json:"organization_type,omitempty"`
    OrganizationDescription *string `json:"organization_description,omitempty"`
    KTP                    *string `json:"ktp,omitempty"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Check if user already exists
    var existingUser models.User
    if err := config.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "error": "User already exists",
        })
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to hash password",
        })
    }

    // Set register status based on role
    registerStatus := "approved"
    if req.Role == "eo" {
        registerStatus = "pending"
    }

    user := models.User{
        Username:                req.Username,
        Name:                    req.Name,
        Email:                   req.Email,
        Password:                string(hashedPassword),
        Role:                    req.Role,
        Organization:            req.Organization,
        OrganizationType:        req.OrganizationType,
        OrganizationDescription: req.OrganizationDescription,
        KTP:                     req.KTP,
        RegisterStatus:          registerStatus,
    }

    if err := config.DB.Create(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create user",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "User registered successfully",
        "user": fiber.Map{
            "user_id":   user.UserID,
            "username":  user.Username,
            "name":      user.Name,
            "email":     user.Email,
            "role":      user.Role,
            "register_status": user.RegisterStatus,
        },
    })
}

func Login(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    var user models.User
    if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }

    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }

    // Check if EO is approved
    if user.Role == "eo" && user.RegisterStatus != "approved" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "EO account not yet approved",
        })
    }

    // Generate JWT
    token, err := utils.GenerateJWT(user.UserID, user.Role)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to generate token",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Login successful",
        "token":   token,
        "user": fiber.Map{
            "user_id":   user.UserID,
            "username":  user.Username,
            "name":      user.Name,
            "email":     user.Email,
            "role":      user.Role,
        },
    })
}