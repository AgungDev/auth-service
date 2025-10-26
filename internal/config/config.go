package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppName           string
	AppPort           string
	DatabaseURL       string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	JWTPrivateKeyPath string
	JWTPublicKeyPath  string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig() // ignore error if .env not found

	cfg := &Config{
		AppName:           viper.GetString("APP_NAME"),
		AppPort:           viper.GetString("APP_PORT"),
		DBHost:            viper.GetString("DB_HOST"),
		DBPort:            viper.GetString("DB_PORT"),
		DBUser:            viper.GetString("DB_USER"),
		DBPassword:        viper.GetString("DB_PASSWORD"),
		DBName:            viper.GetString("DB_NAME"),
		DBSSLMode:         viper.GetString("DB_SSLMODE"),
		JWTPrivateKeyPath: viper.GetString("JWT_PRIVATE_KEY_PATH"),
		JWTPublicKeyPath:  viper.GetString("JWT_PUBLIC_KEY_PATH"),
	}

	cfg.DatabaseURL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	return cfg
}
