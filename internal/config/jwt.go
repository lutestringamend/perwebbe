package config

import (
	"time"

	"github.com/spf13/viper"
)

type JWTConfig struct {
	SecretKey     string        `mapstructure:"JWT_SECRET_KEY"`
	TokenExpiry   time.Duration `mapstructure:"JWT_TOKEN_EXPIRY"`
	RefreshExpiry time.Duration `mapstructure:"JWT_REFRESH_EXPIRY"`
	Issuer        string        `mapstructure:"JWT_ISSUER"`
}

func LoadJWTConfig() (JWTConfig, error) {
	var config JWTConfig

	viper.SetDefault("JWT_TOKEN_EXPIRY", time.Hour*24)
	viper.SetDefault("JWT_REFRESH_EXPIRY", time.Hour*24*7)
	viper.SetDefault("JWT_ISSUER", "chrhndy-perweb")

	config.SecretKey = viper.GetString("JWT_SECRET_KEY")
	config.TokenExpiry = viper.GetDuration("JWT_TOKEN_EXPIRY")
	config.RefreshExpiry = viper.GetDuration("JWT_REFRESH_EXPIRY")
	config.Issuer = viper.GetString("JWT_ISSUER")

	return config, nil
}
