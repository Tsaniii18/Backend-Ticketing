package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "log"
    "ticketing-backend/config"
    "ticketing-backend/controllers"
    "ticketing-backend/middleware"
    "ticketing-backend/models"
)

func main() {
    // Connect to database
    config.ConnectDB()

    // Setup database tanpa foreign key constraints
    setupDatabase()

    app := fiber.New()

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

    // Setup routes
    setupRoutes(app)

    log.Println("🚀 Server running on port 3000")
    log.Println("📊 Database setup completed")
    app.Listen(":3000")
}

func setupDatabase() {
    // Disable foreign key checks
    config.DB.Exec("SET FOREIGN_KEY_CHECKS=0")

    // Auto migrate tanpa foreign key constraints
    err := config.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(
        &models.User{},
        &models.Event{},
        &models.TicketCategory{},
        &models.Report{},
        &models.TransactionHistory{},
        &models.Ticket{},
        &models.Cart{},
    )
    
    if err != nil {
        log.Fatal("❌ Migration failed:", err)
    }

    // Enable foreign key checks kembali
    config.DB.Exec("SET FOREIGN_KEY_CHECKS=1")
    
    log.Println("✅ Database tables created successfully")
}

func setupRoutes(app *fiber.App) {
    // Auth routes
    auth := app.Group("/api/auth")
    auth.Post("/register", controllers.Register)
    auth.Post("/login", controllers.Login)

    // User routes
    user := app.Group("/api/users")
    user.Use(middleware.AuthMiddleware)
    user.Get("/profile", controllers.GetProfile)
    user.Put("/profile", controllers.UpdateProfile)
    user.Get("", middleware.AdminMiddleware, controllers.GetUsers)
    user.Post("/:id/verify", middleware.AdminMiddleware, controllers.VerifyUser)

    // Event routes
    event := app.Group("/api/events")
    event.Get("", controllers.GetEvents)
    event.Get("/:id", controllers.GetEvent)
    
    eventAuth := event.Group("")
    eventAuth.Use(middleware.AuthMiddleware)
    eventAuth.Post("", middleware.EOMiddleware, controllers.CreateEvent)
    eventAuth.Put("/:id", middleware.EOMiddleware, controllers.UpdateEvent)
    eventAuth.Delete("/:id", middleware.EOMiddleware, controllers.DeleteEvent)
    eventAuth.Patch("/:id/verify", middleware.AdminMiddleware, controllers.VerifyEvent)

    // Ticket routes
    ticket := app.Group("/api/tickets")
    ticket.Use(middleware.AuthMiddleware)
    ticket.Post("", controllers.CreateTicket)
    ticket.Get("", controllers.GetTickets)
    ticket.Get("/:id", controllers.GetTicket)
    ticket.Patch("/:id/checkin", controllers.CheckInTicket)

    // Cart routes
    cart := app.Group("/api/cart")
    cart.Use(middleware.AuthMiddleware)
    cart.Post("", controllers.AddToCart)
    cart.Patch("", controllers.UpdateCart)
    cart.Delete("", controllers.DeleteFromCart)
    cart.Post("/checkout", controllers.Checkout)
}