package model

import "time"

// Role represents a role that can be assigned to users within an organization.
type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	OrgID       *uint     `gorm:"index" json:"org_id,omitempty"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	IsSystem    bool      `gorm:"not null;default:false" json:"is_system"`

	// Relationships
	Organization    *Organization    `gorm:"foreignKey:OrgID;references:ID" json:"organization,omitempty"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID;references:ID" json:"role_permissions,omitempty"`
	UserRoles       []UserRole       `gorm:"foreignKey:RoleID;references:ID" json:"user_roles,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}
