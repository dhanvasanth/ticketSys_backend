package services

import (
	"fmt"
	"net/smtp"
	"notifications/internal/config"
	"notifications/internal/repositories"
	"notifications/internal/utils"
	"go.uber.org/zap" 
)

func SendOTP(email string, cfg *config.Config, log *zap.Logger) error {
	code := utils.GenerateOTP(6)
	log.Info("Generated OTP", zap.String("email", email), zap.String("otp", code))

	if err := repositories.SaveOTP(email, code, log); err != nil {
		log.Error("Failed to save OTP", zap.Error(err))
		return err
	}

	auth := smtp.PlainAuth("", cfg.SMTP.Email, cfg.SMTP.Password, "smtp.gmail.com")
	msg := []byte(fmt.Sprintf("Subject: Your OTP\r\n\r\nYour OTP is: %s", code))

	err := smtp.SendMail("smtp.gmail.com:587", auth, cfg.SMTP.Email, []string{email}, msg)
	if err != nil {
		log.Error("Failed to send email", zap.Error(err))
	}
	return err
}

