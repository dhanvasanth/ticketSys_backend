package services

import (
	"fmt"
	"log"
	"user-manager/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db := config.AppConfig.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		db.User, db.Password, db.Host, db.Port, db.Name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// No AutoMigrate() needed since you manually created tables
	log.Println("✅ Connected to database successfully")
}
