package repositories

import (
	"context"

	"auth_service/internal/domain"
	"github.com/google/uuid"
)

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id uuid.UUID, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ClientRepositoryInterface defines the interface for client repository operations
type ClientRepositoryInterface interface {
	Create(ctx context.Context, client *domain.Client) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Client, error)
	FindByClientID(ctx context.Context, clientID string) (*domain.Client, error)
	FindAll(ctx context.Context) ([]domain.Client, error)
	Update(ctx context.Context, id uuid.UUID, client *domain.Client) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// RoleRepositoryInterface defines the interface for role repository operations
type RoleRepositoryInterface interface {
	Create(ctx context.Context, role *domain.Role) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	FindAll(ctx context.Context) ([]domain.Role, error)
	Update(ctx context.Context, id uuid.UUID, role *domain.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PermissionRepositoryInterface defines the interface for permission repository operations
type PermissionRepositoryInterface interface {
	Create(ctx context.Context, permission *domain.Permission) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error)
	FindByName(ctx context.Context, name string) (*domain.Permission, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Permission, error)
	FindAll(ctx context.Context) ([]domain.Permission, error)
	Update(ctx context.Context, id uuid.UUID, permission *domain.Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// TokenRepositoryInterface defines the interface for token repository operations
type TokenRepositoryInterface interface {
	Create(ctx context.Context, token *domain.Token) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Token, error)
	FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*domain.Token, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Token, error)
	Update(ctx context.Context, id uuid.UUID, token *domain.Token) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
