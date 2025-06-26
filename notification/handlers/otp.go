package handlers

import (
	"net/http"
	"notification/config"
	"notification/database"
	"notification/logger"
	"notification/mail"
	"notification/models"
	"notification/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SendOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	otp := utils.GenerateOTP(6)
	expires := time.Now().Add(time.Duration(config.Cfg.OTP.ExpiryMinutes) * time.Minute)

	var existing models.OTPEntry
	err := database.DB.Where("email = ?", req.Email).First(&existing).Error
	if err == nil {
		existing.OTP = otp
		existing.ExpiresAt = expires
		database.DB.Save(&existing)
	} else {
		database.DB.Create(&models.OTPEntry{
			Email:     req.Email,
			OTP:       otp,
			ExpiresAt: expires,
		})
	}

	body := "Your OTP is: " + otp
	if err := mail.SendEmail(req.Email, "Your OTP Code", body); err != nil {
		logger.Log.Error("Failed to send email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	logger.Log.Infof("OTP sent to %s", req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
}

func VerifyOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var entry models.OTPEntry
	err := database.DB.Where("email = ?", req.Email).First(&entry).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No OTP found"})
		return
	}

	if entry.OTP != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect OTP"})
		return
	}

	if time.Now().After(entry.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired"})
		return
	}

	database.DB.Delete(&entry)
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
