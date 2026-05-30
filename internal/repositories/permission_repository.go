package repositories

import (
	"auth_service/internal/domain"
	"errors"

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
func (r *permissionRepository) Create(permission *domain.Permission) error {
	return r.db.Create(permission).Error
}

// FindByID finds a permission by ID
func (r *permissionRepository) FindByID(id uint) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.db.Where("id = ?", id).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// FindByName finds a permission by name
func (r *permissionRepository) FindByName(name string) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.db.Where("name = ?", name).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// FindAll retrieves all permissions
func (r *permissionRepository) FindAll() ([]domain.Permission, error) {
	var permissions []domain.Permission
	if err := r.db.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// Update updates an existing permission
func (r *permissionRepository) Update(id uint, permission *domain.Permission) error {
	return r.db.Model(&domain.Permission{}).Where("id = ?", id).Updates(permission).Error
}

// Delete deletes a permission by ID
func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Permission{}, id).Error
}
