package domain

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ClientID         string
	Name             string
	ClientSecretHash string
	IsConfidential   bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Grants           []ClientGrant       `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
	RedirectURIs     []ClientRedirectURI `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}
