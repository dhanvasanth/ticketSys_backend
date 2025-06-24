package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
    Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
    Host         string `mapstructure:"host"`
    Port         int    `mapstructure:"port"`
    Username     string `mapstructure:"username"`
    Password     string `mapstructure:"password"`
    Database     string `mapstructure:"database"`
    Charset      string `mapstructure:"charset"`
    ParseTime    bool   `mapstructure:"parse_time"`
    MaxIdleConns int    `mapstructure:"max_idle_conns"`
    MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type JWTConfig struct {
    Secret       string `mapstructure:"secret"`
    ExpiresHours int    `mapstructure:"expires_hours"`
}

type LoggingConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")
    viper.AddConfigPath(".")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}