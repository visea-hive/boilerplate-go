package model

import "time"

// EmailVerification represents a token for email verification flow.
type EmailVerification struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt  time.Time  `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UserUUID   string     `gorm:"type:varchar(36);not null;index" json:"user_uuid"`
	TokenHash  string     `gorm:"type:varchar(255);not null" json:"-"`
	ExpiresAt  time.Time  `gorm:"not null" json:"expires_at"`
	VerifiedAt *time.Time `gorm:"" json:"verified_at,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}
