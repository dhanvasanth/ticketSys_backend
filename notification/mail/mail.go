package mail

import (
	"fmt"
	"net/smtp"
	"notification/config"
)

func SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth(
		"",
		config.Cfg.SMTP.Username,
		config.Cfg.SMTP.Password,
		config.Cfg.SMTP.Host,
	)

	
	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			body + "\r\n")

	addr := fmt.Sprintf("%s:%d", config.Cfg.SMTP.Host, config.Cfg.SMTP.Port)

	// ✅ to is passed as recipient
	return smtp.SendMail(addr, auth, config.Cfg.SMTP.From, []string{to}, msg)
}
