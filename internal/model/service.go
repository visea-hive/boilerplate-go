package model

import "time"

// ServiceStatus represents the status of a service.
type ServiceStatus string

const (
	ServiceStatusActive   ServiceStatus = "active"
	ServiceStatusInactive ServiceStatus = "inactive"
)

// Service represents a registered microservice or application.
type Service struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time     `gorm:"autoCreateTime;not null" json:"created_at"`
	CreatedBy   uint          `gorm:"not null" json:"created_by"`
	Name        string        `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description *string       `gorm:"type:text" json:"description,omitempty"`
	URL         string        `gorm:"type:varchar(255);not null" json:"url"`
	Status      ServiceStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// Relationships
	Permissions []Permission `gorm:"foreignKey:SvcID;references:ID" json:"permissions,omitempty"`
}

func (Service) TableName() string {
	return "services"
}
