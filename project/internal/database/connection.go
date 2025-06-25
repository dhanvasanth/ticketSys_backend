// project/internal/database/connection.go
package database

import (
    "fmt"
    "log"
    "project/internal/config"
    "project/internal/models"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.Database,
        cfg.Charset,
        cfg.ParseTime,
    )
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }
    
    // Auto migrate tables
    if err := db.AutoMigrate(
        &models.Role{},
        &models.Permission{},
        &models.User{},
        &models.Ticket{},
        &models.TicketComment{},
        &models.RoleChangeRequest{}, // Add the new model
    ); err != nil {
        return nil, err
    }
    
    // Seed default roles
    seedDefaultRoles(db)
    
    log.Println("Database connected and migrated successfully")
    return db, nil
}

func seedDefaultRoles(db *gorm.DB) {
    roles := []models.Role{
        {
            Name:        "admin",
            Description: "System Administrator",
            Permissions: `["*"]`, // All permissions
        },
        {
            Name:        "agent",
            Description: "Support Agent",
            Permissions: `["ticket.read", "ticket.update", "ticket.assign", "comment.create"]`,
        },
        {
            Name:        "user",
            Description: "Regular User",
            Permissions: `["ticket.create", "ticket.read_own", "comment.create_own"]`,
        },
        {
            Name:        "customer",
            Description: "Customer",
            Permissions: `["ticket.create", "ticket.read_own", "comment.create_own"]`,
        },
    }
    
    for _, role := range roles {
        var existingRole models.Role
        if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
            db.Create(&role)
        }
    }
}