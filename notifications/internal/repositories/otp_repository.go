package repositories

import (
	"notifications/internal/database"
	"notifications/internal/models"
	"go.uber.org/zap" 
	"time"
)

func SaveOTP(email, code string, log *zap.Logger) error {
	otp := models.OTP{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	err := database.DB.Save(&otp).Error
	if err != nil {
		log.Error("Failed to save OTP to DB", zap.Error(err))
	}
	return err
}

func VerifyOTP(email, code string, log *zap.Logger) bool {
	var otp models.OTP
	err := database.DB.Where("email = ? AND code = ?", email, code).First(&otp).Error
	if err != nil {
		log.Warn("OTP not found", zap.String("email", email))
		return false
	}
	if time.Now().After(otp.ExpiresAt) {
		log.Warn("OTP expired", zap.String("email", email))
		database.DB.Delete(&otp)
		return false
	}
	database.DB.Delete(&otp)
	log.Info("OTP verified in DB", zap.String("email", email))
	return true
}
