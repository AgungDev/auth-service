package config

import (
	"fmt"
	"time"

	loggerpkg "auth_service/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB opens a connection to the Postgres database with retry logic.
// This helps when the DB container is still starting while the service tries to connect.
func ConnectDB(dsn, serviceName, serviceVersion, environment string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	maxAttempts := 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: loggerpkg.NewGormLogger(serviceName, serviceVersion, environment, logger.Info),
		})
		if err == nil {
			sqlDB, err2 := db.DB()
			if err2 == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					return db, nil
				} else {
					err = pingErr
				}
			} else {
				err = err2
			}
		}

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return nil, err
}

// GenerateDSN creates a PostgreSQL DSN from config
func GenerateDSN(cfg *Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Name,
		cfg.DBConfig.SSLMode,
	)
}
