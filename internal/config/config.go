package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Driver   string
	SSLMode  string
}

type AppConfig struct {
	Name string
	Port string
}

type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
}

type Config struct {
	AppConfig
	DBConfig
	JWTConfig
}

func (c *Config) setConfig() error {

	// AppConfig
	c.AppConfig.Name = viper.GetString("APP_NAME")
	c.AppConfig.Port = viper.GetString("APP_PORT")

	// DBConfig
	c.DBConfig.Host = viper.GetString("DB_HOST")
	c.DBConfig.Port = viper.GetString("DB_PORT")
	c.DBConfig.User = viper.GetString("DB_USER")
	c.DBConfig.Password = viper.GetString("DB_PASSWORD")
	c.DBConfig.Name = viper.GetString("DB_NAME")
	c.DBConfig.SSLMode = viper.GetString("DB_SSLMODE")
	c.DBConfig.Driver = viper.GetString("DB_DRIVER")
	// JWTConfig
	c.JWTConfig.PrivateKeyPath = viper.GetString("JWT_PRIVATE_KEY_PATH")
	c.JWTConfig.PublicKeyPath = viper.GetString("JWT_PUBLIC_KEY_PATH")

	if c.DBConfig.Host == "" || c.DBConfig.Port == "" || c.DBConfig.User == "" || c.DBConfig.Password == "" || c.AppConfig.Port == "" {
		return fmt.Errorf("Required configuration is missing!")
	}

	// create database URL connection
	// c.DatabaseURL = fmt.Sprintf(
	// 	"postgres://%s:%s@%s:%s/%s?sslmode=%s",	
	// 	c.DBConfig.User,
	// 	c.DBConfig.Password,
	// 	c.DBConfig.Host,
	// 	c.DBConfig.Port,
	// 	c.DBConfig.Name,
	// 	c.DBConfig.SSLMode,
	// )

	return nil
}

func GetConfig() (*Config, error) {
	cfg := &Config{}
	err := cfg.setConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}