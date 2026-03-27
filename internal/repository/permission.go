package repository

import (
	"context"

	"gorm.io/gorm"
)

// PermissionRepository defines the interface for permission-related database operations.
type PermissionRepository interface {
	LoadPermissionsFromDB(ctx context.Context, roleID uint) []string
}

type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new instance of PermissionRepository.
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// LoadPermissionsFromDB queries the database for all valid permission full_keys for a given role.
func (r *permissionRepository) LoadPermissionsFromDB(ctx context.Context, roleID uint) []string {
	var fullKeys []string
	r.db.WithContext(ctx).
		Table("role_permissions").
		Select("permissions.full_key").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Where("permissions.is_deprecated = false").
		Pluck("full_key", &fullKeys)

	return fullKeys
}
