package model

import "time"

// Permission represents a specific permission within a service.
type Permission struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;not null" json:"created_at"`
	CreatedBy    uint       `gorm:"not null" json:"created_by"`
	SvcID        uint       `gorm:"not null;index" json:"svc_id"`
	Key          string     `gorm:"type:varchar(255);not null" json:"key"`
	FullKey      string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"full_key"`
	Description  *string    `gorm:"type:text" json:"description,omitempty"`
	IsDeprecated bool       `gorm:"not null;default:false" json:"is_deprecated"`
	DeprecatedAt *time.Time `gorm:"" json:"deprecated_at,omitempty"`

	// Relationships
	Service Service `gorm:"foreignKey:SvcID;references:ID" json:"service,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}
