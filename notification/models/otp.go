package models

import "time"

type OTPEntry struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"unique"`
	OTP       string
	ExpiresAt time.Time
	CreatedAt time.Time
}
