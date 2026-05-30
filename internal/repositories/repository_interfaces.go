package repositories

import "auth_service/internal/domain"

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(id uint, user *domain.User) error
	Delete(id uint) error
}

// ClientRepositoryInterface defines the interface for client repository operations
type ClientRepositoryInterface interface {
	Create(client *domain.Client) error
	FindByID(id uint) (*domain.Client, error)
	FindByClientID(clientID string) (*domain.Client, error)
	FindAll() ([]domain.Client, error)
	Update(id uint, client *domain.Client) error
	Delete(id uint) error
}

// RoleRepositoryInterface defines the interface for role repository operations
type RoleRepositoryInterface interface {
	Create(role *domain.Role) error
	FindByID(id uint) (*domain.Role, error)
	FindByName(name string) (*domain.Role, error)
	FindAll() ([]domain.Role, error)
	Update(id uint, role *domain.Role) error
	Delete(id uint) error
}

// PermissionRepositoryInterface defines the interface for permission repository operations
type PermissionRepositoryInterface interface {
	Create(permission *domain.Permission) error
	FindByID(id uint) (*domain.Permission, error)
	FindByName(name string) (*domain.Permission, error)
	FindAll() ([]domain.Permission, error)
	Update(id uint, permission *domain.Permission) error
	Delete(id uint) error
}

// TokenRepositoryInterface defines the interface for token repository operations
type TokenRepositoryInterface interface {
	Create(token *domain.Token) error
	FindByID(id uint) (*domain.Token, error)
	FindByAccessToken(accessToken string) (*domain.Token, error)
	FindByUserID(userID uint) (*domain.Token, error)
	Update(id uint, token *domain.Token) error
	Delete(id uint) error
}
