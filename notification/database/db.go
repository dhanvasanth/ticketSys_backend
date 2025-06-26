package database

import (
	"notification/config"
	"notification/logger"
	"notification/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := config.GetDSN()
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Database connection failed:", err)
	}
	DB.AutoMigrate(&models.OTPEntry{})
}
