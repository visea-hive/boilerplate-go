package model

import "time"

// UserStatus represents the status of a user.
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusPending   UserStatus = "pending"
	UserStatusSuspended UserStatus = "suspended"
)

// User represents an authenticated user in the system.
type User struct {
	BaseModel

	Email           string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	EmailVerifiedAt *time.Time `gorm:"" json:"email_verified_at,omitempty"`
	Password        *string    `gorm:"type:varchar(255)" json:"-"`
	DisplayName     *string    `gorm:"type:varchar(255)" json:"display_name,omitempty"`
	AvatarURL       *string    `gorm:"type:text" json:"avatar_url,omitempty"`
	Status          UserStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	MFAEnabled      bool       `gorm:"not null;default:false" json:"mfa_enabled"`
	LastLoginAt     *time.Time `gorm:"" json:"last_login_at,omitempty"`
	OrganizationID  *uint      `gorm:"index" json:"organization_id,omitempty"`
	Lang            *string    `gorm:"type:varchar(10)" json:"lang,omitempty"`

	// Relationships
	Organization *Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"organization,omitempty"`
	Memberships  []Membership  `gorm:"foreignKey:UserUUID;references:UUID" json:"memberships,omitempty"`
}

// TableName specifies the table name for the User model.
func (User) TableName() string {
	return "users"
}

// UserDTO represents the API response for a user.
type UserDTO struct {
	ID              uint            `json:"id"`
	Email           string          `json:"email"`
	EmailVerifiedAt *time.Time      `json:"email_verified_at"`
	DisplayName     *string         `json:"display_name"`
	AvatarURL       *string         `json:"avatar_url"`
	Status          string          `json:"status"`
	MFAEnabled      bool            `json:"mfa_enabled"`
	LastLoginAt     *time.Time      `json:"last_login_at"`
	OrganizationID  *uint           `json:"organization_id"`
	Organization    OrganizationDTO `json:"organization"`
	Lang            *string         `json:"lang"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
}

func MapToUserDTO(user *User) *UserDTO {
	if user == nil {
		return nil
	}

	var organizationDTO OrganizationDTO
	if user.Organization != nil {
		organizationDTO = MapToOrganizationDTO(user.Organization)
	}

	return &UserDTO{
		ID:              user.ID,
		Email:           user.Email,
		EmailVerifiedAt: user.EmailVerifiedAt,
		DisplayName:     user.DisplayName,
		AvatarURL:       user.AvatarURL,
		Status:          string(user.Status),
		MFAEnabled:      user.MFAEnabled,
		LastLoginAt:     user.LastLoginAt,
		OrganizationID:  user.OrganizationID,
		Organization:    organizationDTO,
		Lang:            user.Lang,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func MapToUserDTOs(users []User) []UserDTO {
	var dtos []UserDTO
	for _, user := range users {
		dto := MapToUserDTO(&user)
		dtos = append(dtos, *dto)
	}

	return dtos
}
