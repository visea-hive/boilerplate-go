package model

import "time"

// RolePermission represents a many-to-many link between roles and permissions.
type RolePermission struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt    time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	CreatedBy    uint      `gorm:"not null" json:"created_by"`
	RoleID       uint      `gorm:"not null;index" json:"role_id"`
	PermissionID uint      `gorm:"not null;index" json:"permission_id"`

	// Relationships
	Role       Role       `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID;references:ID" json:"permission,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
