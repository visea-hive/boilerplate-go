package model

import "time"

// RevokeReason represents the reason a session was revoked.
type RevokeReason string

const (
	RevokeReasonLogout     RevokeReason = "logout"
	RevokeReasonForced     RevokeReason = "forced"
	RevokeReasonSuspicious RevokeReason = "suspicious"
	RevokeReasonExpired    RevokeReason = "expired"
)

// Session represents an authenticated user session.
type Session struct {
	ID           uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt    time.Time     `gorm:"autoCreateTime;not null" json:"created_at"`
	UserUUID     string        `gorm:"type:varchar(36);not null;index" json:"user_uuid"`
	OrgID        *uint         `gorm:"index" json:"org_id,omitempty"`
	DeviceID     *uint         `gorm:"index" json:"device_id,omitempty"`
	IPAddress    *string       `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent    *string       `gorm:"type:text" json:"user_agent,omitempty"`
	LastActiveAt *time.Time    `gorm:"" json:"last_active_at,omitempty"`
	ExpiresAt    time.Time     `gorm:"not null" json:"expires_at"`
	RevokedAt    *time.Time    `gorm:"" json:"revoked_at,omitempty"`
	RevokeReason *RevokeReason `gorm:"type:varchar(20)" json:"revoke_reason,omitempty"`

	// Relationships
	User          User           `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
	Organization  *Organization  `gorm:"foreignKey:OrgID;references:ID" json:"organization,omitempty"`
	Device        *Device        `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:SessionID;references:ID" json:"refresh_tokens,omitempty"`
}

func (Session) TableName() string {
	return "sessions"
}
