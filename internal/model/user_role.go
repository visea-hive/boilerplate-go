package model

import "time"

// UserRoleStatus represents the status of a user role assignment.
type UserRoleStatus string

const (
	UserRoleStatusActive    UserRoleStatus = "active"
	UserRoleStatusSuspended UserRoleStatus = "suspended"
)

// UserRole represents a user's role assignment scoped to an organization.
type UserRole struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	CreatedBy uint           `gorm:"not null" json:"created_by"`
	UserUUID  string         `gorm:"type:varchar(36);not null;index" json:"user_uuid"`
	RoleID    uint           `gorm:"not null;index" json:"role_id"`
	OrgID     uint           `gorm:"not null;index" json:"org_id"`
	Status    UserRoleStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	ExpiresAt *time.Time     `gorm:"" json:"expires_at,omitempty"`

	// Relationships
	User         User         `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
	Role         Role         `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
	Organization Organization `gorm:"foreignKey:OrgID;references:ID" json:"organization,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
