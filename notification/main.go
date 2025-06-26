package main


import (
	"github.com/gin-gonic/gin"
	"notification/config"
	"notification/database"
	"notification/handlers"
	"notification/logger"
)

func main() {
	config.LoadConfig()
	logger.InitLogger()
	database.InitDB()

	r := gin.Default()

	r.POST("/send-otp", handlers.SendOTP)
	r.POST("/verify-otp", handlers.VerifyOTP)

	logger.Log.Info("Server running on http://localhost:8081")
	r.Run(":8081")
}

