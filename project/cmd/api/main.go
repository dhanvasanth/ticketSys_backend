package main

import (
    "fmt"
    "log"
    "project/internal/config"
    "project/internal/database"
    "project/internal/handlers"
    "project/internal/middleware"
    "project/internal/repositories"
    "project/internal/services"
    
    "github.com/gin-gonic/gin"
)

func main() {
    // Load configurationproject/
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Connect to database
    db, err := database.Connect(&cfg.Database)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Initialize repositories
    userRepo := repositories.NewUserRepository(db)
    ticketRepo := repositories.NewTicketRepository(db)
    
    // Initialize services
    authService := services.NewAuthService(userRepo, &cfg.JWT)
    userService := services.NewUserService(userRepo)
    ticketService := services.NewTicketService(ticketRepo)
    
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService)
    userHandler := handlers.NewUserHandler(userService)
    ticketHandler := handlers.NewTicketHandler(ticketService)
    
    // Setup router
    gin.SetMode(cfg.Server.Mode)
    r := gin.Default()
    
    // Public routes
    auth := r.Group("/api/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
    }
    
    // Protected routes
    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware(&cfg.JWT))
    {
        // User routes
        users := api.Group("/users")
        {
            users.GET("/profile", userHandler.GetProfile)
            users.PUT("/profile", userHandler.UpdateProfile)
        }
        
        // Ticket routes with permissions
        tickets := api.Group("/tickets")
        {
            tickets.GET("", ticketHandler.GetTickets) // Users see only their tickets
            tickets.GET("/:id", ticketHandler.GetTicket)
            tickets.POST("", middleware.RequirePermission("ticket.create"), ticketHandler.CreateTicket)
            tickets.PUT("/:id", middleware.RequirePermission("ticket.update"), ticketHandler.UpdateTicket)
            tickets.POST("/:id/comments", middleware.RequirePermission("comment.create"), ticketHandler.AddComment)
        }
        
        // Admin routes
        admin := api.Group("/admin")
        admin.Use(middleware.RequirePermission("*")) // Only admins
        {
            admin.GET("/users", userHandler.GetAllUsers)
            admin.GET("/tickets", ticketHandler.GetAllTickets) // All tickets
            admin.POST("/roles", userHandler.CreateRole)
            admin.GET("/roles", userHandler.GetRoles)
        }
    }
    
    // Start server
    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("Server starting on %s", addr)
    log.Fatal(r.Run(addr))
}