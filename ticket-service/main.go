package main

import(
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)



func main() {
	// Initialize the global logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	dsn := "mariadb:mariadb@tcp(localhost:3306)/tickets?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Fatalf("Failed to connect to MariaDB: %v", err)
	}
}
