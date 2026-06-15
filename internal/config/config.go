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
	Name        string
	Port        string
	Version     string
	Environment string
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

func LoadConfig() error {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}

func (c *Config) setConfig() error {
	if err := LoadConfig(); err != nil {
		return err
	}

	// AppConfig
	c.AppConfig.Name = viper.GetString("APP_NAME")
	c.AppConfig.Port = viper.GetString("APP_PORT")
	c.AppConfig.Version = viper.GetString("APP_VERSION")
	c.AppConfig.Environment = viper.GetString("APP_ENV")
	if c.AppConfig.Environment == "" {
		c.AppConfig.Environment = "development"
	}

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

	if c.AppConfig.Name == "" || c.AppConfig.Port == "" || c.AppConfig.Version == "" {
		return fmt.Errorf("required application configuration is missing")
	}

	if c.DBConfig.Host == "" || c.DBConfig.Port == "" || c.DBConfig.User == "" || c.DBConfig.Password == "" {
		return fmt.Errorf("required database configuration is missing")
	}

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
