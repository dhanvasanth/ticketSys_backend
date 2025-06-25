package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"db"`

	SMTP struct {
		Email    string `mapstructure:"email"`
		Password string `mapstructure:"password"`
	} `mapstructure:"smtp"`
}

func LoadConfig(path string) *Config {
	var cfg Config
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config:", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Unable to decode config:", err)
	}
	return &cfg
}
