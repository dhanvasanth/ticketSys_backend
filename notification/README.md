# 🔐 Go OTP API Microservice

This is a Go-based OTP microservice that:

- Sends OTP to user via email using SMTP (e.g., Gmail)
- Verifies OTP with expiration
- Uses GORM with MariaDB for database
- Uses Zap for structured logging
- Loads config from a `config.yaml` file

---

## 📁 Project Structure

otp-api/
├── config.yaml # Configuration for DB, SMTP, OTP
├── main.go # Entry point
├── go.mod / go.sum # Go modules
├── config/ # Loads config.yaml
├── database/ # GORM DB init
├── handlers/ # Gin handlers
├── logger/ # Zap logging
├── mail/ # SMTP email sender
├── models/ # OTPEntry GORM model
├── utils/ # OTP generator

yaml
Copy
Edit

---

## ⚙️ Prerequisites

- Go 1.18+
- MariaDB running (local or Docker)
- Gmail account with App Password if using Gmail SMTP
- `config.yaml` in project root

---

## 🔧 Configuration

### `config.yaml`

```yaml
smtp:
  host: smtp.gmail.com
  port: 587
  username: your_email@gmail.com
  password: your_app_password
  from: your_email@gmail.com

otp:
  expiry_minutes: 5

database:
  host: localhost
  port: 3306
  user: mariadb
  password: pass
  name: otp
🔐 Use a Gmail App Password if 2FA is enabled

🐳 Running MariaDB in Docker
bash
Copy
Edit

🚀 Running the API
bash
Copy
Edit
go mod tidy
go run main.go
Logs will print to console via Zap logger

📮 API Endpoints
🔸 Send OTP
bash
Copy
Edit
POST /send-otp

http://localhost:8081/send-otp
Body (JSON):{
  "email": "user@example.com"
}

🔸 Verify OTP
bash
Copy
Edit
POST http://localhost:8081/verify-otp
Body (JSON):

json
Copy
Edit
{
  "email": "user@example.com",
  "otp": "123456"
}