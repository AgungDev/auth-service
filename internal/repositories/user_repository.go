package repositories

import (
	"auth_service/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindAll retrieves all users
func (r *userRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, id uuid.UUID, user *domain.User) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(user).Error
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}

// AssignRoles assigns roles to a user
func (r *userRepository) AssignRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	for _, roleID := range roleIDs {
		if err := r.db.WithContext(ctx).Exec(
			"INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
			userID,
			roleID,
		).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetRolesByUserID retrieves roles assigned to a user
func (r *userRepository) GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.WithContext(ctx).Raw(
		`SELECT r.*
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`,
		userID,
	).Scan(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// RemoveRole removes a role assignment from a user
func (r *userRepository) RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM user_roles WHERE user_id = ? AND role_id = ?",
		userID,
		roleID,
	).Error
}

// AssignPermissions assigns direct permissions to a user
func (r *userRepository) AssignPermissions(ctx context.Context, userID uuid.UUID, permissionIDs []uuid.UUID) error {
	for _, permissionID := range permissionIDs {
		if err := r.db.WithContext(ctx).Exec(
			"INSERT INTO user_permissions (user_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
			userID,
			permissionID,
		).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetPermissions retrieves direct permissions assigned to a user
func (r *userRepository) GetPermissions(ctx context.Context, userID uuid.UUID) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).Raw(
		`SELECT p.*
		FROM permissions p
		JOIN user_permissions up ON up.permission_id = p.id
		WHERE up.user_id = ?`,
		userID,
	).Scan(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// RemovePermission removes a direct permission from a user
func (r *userRepository) RemovePermission(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM user_permissions WHERE user_id = ? AND permission_id = ?",
		userID,
		permissionID,
	).Error
}
