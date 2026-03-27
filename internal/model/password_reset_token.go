package model

import "time"

// PasswordResetToken represents a token for password reset flow.
type PasswordResetToken struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UserUUID  string     `gorm:"type:varchar(36);not null;index" json:"user_uuid"`
	TokenHash string     `gorm:"type:varchar(255);not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `gorm:"" json:"used_at,omitempty"`
	IPAddress *string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
