package model

import "time"

// RefreshToken represents a refresh token for session renewal.
type RefreshToken struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	IssuedAt     time.Time  `gorm:"autoCreateTime;not null" json:"issued_at"`
	SessionID    uint       `gorm:"not null;index" json:"session_id"`
	TokenHash    string     `gorm:"type:varchar(255);not null" json:"-"`
	Family       string     `gorm:"type:varchar(36);not null" json:"family"`
	ReplacedByID *uint      `gorm:"index" json:"replaced_by_id,omitempty"`
	ExpiresAt    time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt       *time.Time `gorm:"" json:"used_at,omitempty"`
	RevokedAt    *time.Time `gorm:"" json:"revoked_at,omitempty"`

	// Relationships
	Session    Session       `gorm:"foreignKey:SessionID;references:ID" json:"session,omitempty"`
	ReplacedBy *RefreshToken `gorm:"foreignKey:ReplacedByID;references:ID" json:"replaced_by,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
