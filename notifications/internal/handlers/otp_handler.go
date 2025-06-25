package handlers

import (
	"net/http"
	"notifications/internal/config"
	"notifications/internal/repositories"
	"notifications/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OTPRequest struct {
	Email string `json:"email"`
}

type OTPVerify struct {
	Email string `json:"email"`
	Code  string `json:"otp"`
}

func SendOTPHandler(cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn("Invalid email input", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}
		if err := services.SendOTP(req.Email, cfg, log); err != nil {
			log.Error("Failed to send OTP", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP send failed"})
			return
		}
		log.Info("OTP sent", zap.String("email", req.Email))
		c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
	}
}

func VerifyOTPHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OTPVerify
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn("Invalid verification request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		if !repositories.VerifyOTP(req.Email, req.Code, log) {
			log.Warn("OTP verification failed", zap.String("email", req.Email))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
			return
		}
		log.Info("OTP verified", zap.String("email", req.Email))
		c.JSON(http.StatusOK, gin.H{"message": "OTP verified"})
	}
}
