package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SMTP struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		From     string `yaml:"from"`
	} `yaml:"smtp"`
	OTP struct {
		ExpiryMinutes int `yaml:"expiry_minutes"`
	} `yaml:"otp"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

var Cfg Config

func LoadConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic("config.yaml not found")
	}
	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		panic("invalid config.yaml")
	}
}

func GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		Cfg.Database.User,
		Cfg.Database.Password,
		Cfg.Database.Host,
		Cfg.Database.Port,
		Cfg.Database.Name,
	)
}
