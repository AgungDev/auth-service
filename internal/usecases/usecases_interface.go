package usecases

import (
	"context"

	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"github.com/google/uuid"
)

// UserUsecaseInterface defines the interface for user usecase operations
type UserUsecaseInterface interface {
	Register(ctx context.Context, req dto.UserRegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req dto.UserLoginRequest) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UserUpdateRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// AuthUsecaseInterface defines the interface for auth-specific business logic
type AuthUsecaseInterface interface {
	Register(ctx context.Context, req dto.UserRegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req dto.UserLoginRequest) (*dto.AuthLoginResponse, error)
	Refresh(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	Introspect(ctx context.Context, token string) (*dto.IntrospectResponse, error)
	GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]string, error)
	Authorize(ctx context.Context, subject string, roles []interface{}, req dto.AuthorizeRequest) (*dto.AuthorizeResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
}

// ClientUsecaseInterface defines the interface for client usecase operations
type ClientUsecaseInterface interface {
	Create(ctx context.Context, req dto.ClientRequest) (*domain.Client, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error)
	GetByClientID(ctx context.Context, clientID string) (*domain.Client, error)
	GetAll(ctx context.Context) ([]domain.Client, error)
	Update(ctx context.Context, id uuid.UUID, req dto.ClientRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PermissionUsecaseInterface defines the interface for permission usecase operations
type PermissionUsecaseInterface interface {
	Create(ctx context.Context, req dto.PermissionRequest) (*domain.Permission, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error)
	GetByName(ctx context.Context, name string) (*domain.Permission, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Permission, error)
	GetAll(ctx context.Context) ([]domain.Permission, error)
	Update(ctx context.Context, id uuid.UUID, req dto.PermissionRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// RoleUsecaseInterface defines the interface for role usecase operations
type RoleUsecaseInterface interface {
	Create(ctx context.Context, req dto.RoleRequest) (*domain.Role, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	GetByName(ctx context.Context, name string) (*domain.Role, error)
	GetAll(ctx context.Context) ([]domain.Role, error)
	Update(ctx context.Context, id uuid.UUID, req dto.RoleRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// TokenUsecaseInterface defines the interface for token usecase operations
type TokenUsecaseInterface interface {
	Create(ctx context.Context, token *domain.Token) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Token, error)
	GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*domain.Token, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Token, error)
	Update(ctx context.Context, id uuid.UUID, token *domain.Token) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	IsTokenExpired(token *domain.Token) bool
}
