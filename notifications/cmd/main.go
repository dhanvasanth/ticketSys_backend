package main

import (
	"github.com/gin-gonic/gin"
	"notifications/internal/config"
	"notifications/internal/database"
	"notifications/internal/handlers"
	"notifications/internal/logger"
)

func main() {
	log := logger.NewZapLogger()
	cfg := config.LoadConfig("config.yaml")
	database.Connect(cfg, log)

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/send-otp", handlers.SendOTPHandler(cfg, log))
	r.POST("/verify-otp", handlers.VerifyOTPHandler(log))

	log.Info("Notification service running on port 8081")
	r.Run(":8081")
}
