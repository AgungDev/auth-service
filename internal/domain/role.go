package domain

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code        string
	Name        string
	Description string
}
