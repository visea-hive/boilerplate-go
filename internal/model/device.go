package model

import "time"

// DevicePlatform represents the platform of a device.
type DevicePlatform string

const (
	DevicePlatformIOS     DevicePlatform = "ios"
	DevicePlatformAndroid DevicePlatform = "android"
	DevicePlatformWeb     DevicePlatform = "web"
)

// Device represents a user's registered device.
type Device struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt         time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	UserUUID          string         `gorm:"type:varchar(36);not null;index" json:"user_uuid"`
	DeviceFingerprint string         `gorm:"type:varchar(255);not null" json:"device_fingerprint"`
	Platform          DevicePlatform `gorm:"type:varchar(20);not null" json:"platform"`
	PushToken         *string        `gorm:"type:text" json:"push_token,omitempty"`
	UserAgent         *string        `gorm:"type:text" json:"user_agent,omitempty"`
	TrustedAt         *time.Time     `gorm:"" json:"trusted_at,omitempty"`
	LastSeenAt        *time.Time     `gorm:"" json:"last_seen_at,omitempty"`
	RevokedAt         *time.Time     `gorm:"" json:"revoked_at,omitempty"`

	// Relationships
	User     User      `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
	Sessions []Session `gorm:"foreignKey:DeviceID;references:ID" json:"sessions,omitempty"`
}

func (Device) TableName() string {
	return "devices"
}
