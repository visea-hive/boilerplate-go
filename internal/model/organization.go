package model

import (
	"time"

	"gorm.io/gorm"
)

// OrganizationStatus represents the status of an organization.
type OrganizationStatus string

const (
	StatusActive    OrganizationStatus = "active"
	StatusSuspended OrganizationStatus = "suspended"
)

// Organization represents a hierarchical organizational entity.
type Organization struct {
	BaseModel

	Name        string             `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string             `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description *string            `gorm:"type:text" json:"description,omitempty"`
	Level       int                `gorm:"not null" json:"level"`
	ParentID    *uint              `gorm:"index" json:"parent_id,omitempty"`
	Status      OrganizationStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	SSOEnforced bool               `gorm:"not null;default:false" json:"sso_enforced"`
	MFAEnforced bool               `gorm:"not null;default:false" json:"mfa_enforced"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy *string        `gorm:"type:varchar(36)" json:"deleted_by,omitempty"`

	// Self-referential relationships
	Parent   *Organization  `gorm:"foreignKey:ParentID;references:ID" json:"parent,omitempty"`
	Children []Organization `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`

	// Reverse relationships
	Users       []User       `gorm:"foreignKey:OrganizationID;references:ID" json:"users,omitempty"`
	Memberships []Membership `gorm:"foreignKey:OrgID;references:ID" json:"memberships,omitempty"`
}

// TableName specifies the table name for the Organization model.
func (Organization) TableName() string {
	return "organizations"
}

// OrganizationDTO represents the API response for an organization.
type OrganizationDTO struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description *string          `json:"description"`
	Level       int              `json:"level"`
	ParentID    *uint            `json:"parent_id"`
	Parent      *OrganizationDTO `json:"parent"`
	Status      string           `json:"status"`
	SSOEnforced bool             `json:"sso_enforced"`
	MFAEnforced bool             `json:"mfa_enforced"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   *time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at,omitempty"`
	DeletedBy   *string          `json:"deleted_by,omitempty"`
}

func MapToOrganizationDTO(org *Organization) OrganizationDTO {
	if org == nil {
		return OrganizationDTO{}
	}

	var parentDTO *OrganizationDTO
	if org.Parent != nil {
		pDTO := MapToOrganizationDTO(org.Parent)
		parentDTO = &pDTO
	}

	return OrganizationDTO{
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		Description: org.Description,
		Level:       org.Level,
		ParentID:    org.ParentID,
		Parent:      parentDTO,
		Status:      string(org.Status),
		SSOEnforced: org.SSOEnforced,
		MFAEnforced: org.MFAEnforced,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
		DeletedAt:   org.DeletedAt,
		DeletedBy:   org.DeletedBy,
	}
}

func MapToOrganizationDTOs(orgs []Organization) []OrganizationDTO {
	var dtos []OrganizationDTO
	for _, org := range orgs {
		dto := MapToOrganizationDTO(&org)
		dtos = append(dtos, dto)
	}

	return dtos
}
