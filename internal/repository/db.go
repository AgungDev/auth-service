package repository

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB opens a connection to the Postgres database with simple retry logic.
// This helps when the DB container is still starting while the service tries to connect.
func NewDB(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	maxAttempts := 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
