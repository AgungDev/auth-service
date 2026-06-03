package domain

import "github.com/google/uuid"

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code        string
	Name        string
	Description string
}
