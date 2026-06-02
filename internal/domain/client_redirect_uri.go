package domain

import (
    "time"

    "github.com/google/uuid"
)

type ClientRedirectURI struct {
    ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    ClientID   uuid.UUID `gorm:"type:uuid;not null;index"`
    RedirectURI string
    CreatedAt  time.Time
}
