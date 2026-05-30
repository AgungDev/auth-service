package usecases

import (
	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
)

// UserUsecaseInterface defines the interface for user usecase operations
type UserUsecaseInterface interface {
	Register(req dto.UserRegisterRequest) (*domain.User, error)
	Login(req dto.UserLoginRequest) (*domain.User, error)
	GetByID(id uint) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetAll() ([]domain.User, error)
	Update(id uint, req dto.UserUpdateRequest) error
	Delete(id uint) error
}

// ClientUsecaseInterface defines the interface for client usecase operations
type ClientUsecaseInterface interface {
	Create(req dto.ClientRequest) (*domain.Client, error)
	GetByID(id uint) (*domain.Client, error)
	GetByClientID(clientID string) (*domain.Client, error)
	GetAll() ([]domain.Client, error)
	Update(id uint, req dto.ClientRequest) error
	Delete(id uint) error
}

// PermissionUsecaseInterface defines the interface for permission usecase operations
type PermissionUsecaseInterface interface {
	Create(req dto.PermissionRequest) (*domain.Permission, error)
	GetByID(id uint) (*domain.Permission, error)
	GetByName(name string) (*domain.Permission, error)
	GetAll() ([]domain.Permission, error)
	Update(id uint, req dto.PermissionRequest) error
	Delete(id uint) error
}

// RoleUsecaseInterface defines the interface for role usecase operations
type RoleUsecaseInterface interface {
	Create(req dto.RoleRequest) (*domain.Role, error)
	GetByID(id uint) (*domain.Role, error)
	GetByName(name string) (*domain.Role, error)
	GetAll() ([]domain.Role, error)
	Update(id uint, req dto.RoleRequest) error
	Delete(id uint) error
}

// TokenUsecaseInterface defines the interface for token usecase operations
type TokenUsecaseInterface interface {
	Create(token *domain.Token) error
	GetByID(id uint) (*domain.Token, error)
	GetByAccessToken(accessToken string) (*domain.Token, error)
	GetByUserID(userID uint) (*domain.Token, error)
	Update(id uint, token *domain.Token) error
	Delete(id uint) error
	IsTokenExpired(token *domain.Token) bool
}
