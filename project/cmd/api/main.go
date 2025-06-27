// project/cmd/api/main.go
package main

import (
    "fmt"
    "project/internal/config"
    "project/internal/database"
    "project/internal/handlers"
    "project/internal/middleware"
    "project/internal/repositories"
    "project/internal/services"
    "github.com/gin-contrib/cors"
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    logger, err := zap.NewProduction()
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize logger: %v", err))
    }
    defer logger.Sync()
    
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        logger.Fatal("Failed to load config", zap.Error(err))
    }
    
    // Connect to database
    db, err := database.Connect(&cfg.Database)
    if err != nil {
        logger.Fatal("Failed to connect to database", zap.Error(err))
    }
    
    // Initialize repositories
    userRepo := repositories.NewUserRepository(db)
    ticketRepo := repositories.NewTicketRepository(db)
    roleChangeRepo := repositories.NewRoleChangeRequestRepository(db)
    
    // Initialize services
    authService := services.NewAuthService(userRepo, &cfg.JWT)
    userService := services.NewUserService(userRepo)
    ticketService := services.NewTicketService(ticketRepo)
    roleChangeService := services.NewRoleChangeRequestService(roleChangeRepo, userRepo)
    
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService)
    userHandler := handlers.NewUserHandler(userService)
    ticketHandler := handlers.NewTicketHandler(ticketService)
    roleChangeHandler := handlers.NewRoleChangeRequestHandler(roleChangeService)
    
    // Setup router
    gin.SetMode(cfg.Server.Mode)
    r := gin.Default()
    r.Use(cors.Default())
    
    // Public routes
    auth := r.Group("/api/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
    }
    
    // Protected routes
    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware(&cfg.JWT, logger))
    {
        // User routes
        users := api.Group("/users")
        {
            users.GET("/profile", userHandler.GetProfile)
            users.PUT("/profile", userHandler.UpdateProfile)
        }
        
        // Role change request routes
        roleRequests := api.Group("/role-requests")
        {
            // Any authenticated user can create and view their own requests
            roleRequests.POST("", roleChangeHandler.CreateRequest)
            roleRequests.GET("/my", roleChangeHandler.GetMyRequests)
            roleRequests.GET("/:id", roleChangeHandler.GetRequest)
        }
        
        // Ticket routes with permissions
        tickets := api.Group("/tickets")
        {
            tickets.GET("", ticketHandler.GetTickets) // Users see only their tickets
            tickets.GET("/:id", ticketHandler.GetTicket)
            tickets.POST("", middleware.RequirePermission("ticket.create", logger), ticketHandler.CreateTicket)
            tickets.PUT("/:id", middleware.RequirePermission("ticket.update", logger), ticketHandler.UpdateTicket)
            tickets.POST("/:id/comments", middleware.RequirePermission("comment.create", logger), ticketHandler.AddComment)
        }
        
        // Admin routes
        admin := api.Group("/admin")
        admin.Use(middleware.RequirePermission("*", logger)) // Only admins
        {
            admin.GET("/users", userHandler.GetAllUsers)
            admin.GET("/tickets", ticketHandler.GetAllTickets) // All tickets
            admin.POST("/roles", userHandler.CreateRole)
            admin.GET("/roles", userHandler.GetRoles)
            
            // Role change request management
            admin.GET("/role-requests", roleChangeHandler.GetAllRequests) // ?status=pending for pending only
            admin.PUT("/role-requests/:id/process", roleChangeHandler.ProcessRequest)
        }
    }
    
    // Start server
    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    logger.Info("Server starting", 
        zap.String("address", addr),
        zap.String("mode", cfg.Server.Mode))
    
    if err := r.Run(addr); err != nil {
        logger.Fatal("Failed to start server", zap.Error(err))
    }
}