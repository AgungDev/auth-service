package repositories

import (
	"auth_service/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new permission repository instance
func NewPermissionRepository(db *gorm.DB) PermissionRepositoryInterface {
	return &permissionRepository{db: db}
}

// Create inserts a new permission into the database
func (r *permissionRepository) Create(ctx context.Context, permission *domain.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

// FindByID finds a permission by ID
func (r *permissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// FindByName finds a permission by name
func (r *permissionRepository) FindByName(ctx context.Context, name string) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// FindAll retrieves all permissions
func (r *permissionRepository) FindAll(ctx context.Context) ([]domain.Permission, error) {
	var permissions []domain.Permission
	if err := r.db.WithContext(ctx).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// FindByUserID retrieves permissions assigned to the user's roles
func (r *permissionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).Raw(
		`SELECT DISTINCT p.*
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ?`, userID).Scan(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// Update updates an existing permission
func (r *permissionRepository) Update(ctx context.Context, id uuid.UUID, permission *domain.Permission) error {
	return r.db.WithContext(ctx).Model(&domain.Permission{}).Where("id = ?", id).Updates(permission).Error
}

// Delete deletes a permission by ID
func (r *permissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Permission{}, id).Error
}
