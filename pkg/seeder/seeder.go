package seeder

import (
	"auth_service/internal/domain"
	"auth_service/pkg"
	loggerpkg "auth_service/pkg/logger"
	"context"
	"gorm.io/gorm"
)

// SeedDummyData inserts initial dummy data if the database is empty.
func SeedDummyData(db *gorm.DB) error {
	ctx := context.WithValue(context.Background(), loggerpkg.ContextKeyLogSource, "startup")
	db = db.WithContext(ctx)
	var count int64

	if err := db.Model(&domain.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		passwordHash, err := pkg.HashPassword("Password123!")
		if err != nil {
			return err
		}
		user := &domain.User{
			Username:     "admin",
			Email:        "admin@example.com",
			PasswordHash: passwordHash,
			FullName:     "Admin User",
			Status:       "active",
		}
		if err := db.Create(user).Error; err != nil {
			return err
		}
	}

	if err := db.Model(&domain.Client{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		secretHash, err := pkg.HashPassword("secret123")
		if err != nil {
			return err
		}
		client := &domain.Client{
			ClientID:         "default-client",
			Name:             "Default Client",
			ClientSecretHash: secretHash,
			IsConfidential:   true,
		}
		if err := db.Create(client).Error; err != nil {
			return err
		}

		// insert normalized redirect URIs
		redirect := &domain.ClientRedirectURI{
			ClientID:   client.ID,
			RedirectURI: "http://localhost/callback",
		}
		if err := db.Create(redirect).Error; err != nil {
			return err
		}

		// insert grants
		grants := []string{"password", "refresh_token"}
		for _, g := range grants {
			cg := &domain.ClientGrant{
				ClientID:  client.ID,
				GrantType: g,
			}
			if err := db.Create(cg).Error; err != nil {
				return err
			}
		}
	}

	if err := db.Model(&domain.Role{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		role := &domain.Role{
			Code:        "admin",
			Name:        "admin",
			Description: "Administrator role",
		}
		if err := db.Create(role).Error; err != nil {
			return err
		}
	}

	if err := db.Model(&domain.Permission{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		permission := &domain.Permission{
			Code:        "all:access",
			Name:        "all:access",
			Description: "Full access permission",
		}
		if err := db.Create(permission).Error; err != nil {
			return err
		}
	}

	return nil
}
