package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenerateUUIDv7 generates a new UUID v7 string.
// This function is reusable across all models.
func GenerateUUIDv7() string {
	return uuid.Must(uuid.NewV7()).String()
}

// BaseModel contains shared fields for all models.
type BaseModel struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      string     `gorm:"type:varchar(36);uniqueIndex;not null" json:"uuid"`
	CreatedAt time.Time  `gorm:"autoCreateTime;not null" json:"created_at"`
	CreatedBy string     `gorm:"type:varchar(36);not null" json:"created_by"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy *string    `gorm:"type:varchar(36)" json:"updated_by,omitempty"`
}

// BeforeCreate hook auto-generates UUID v7 for the UUID field.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.UUID == "" {
		b.UUID = GenerateUUIDv7()
	}
	return nil
}
