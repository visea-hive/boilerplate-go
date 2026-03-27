package model

import "time"

// OTPChannel represents the channel used for OTP delivery.
type OTPChannel string

const (
	OTPChannelEmail OTPChannel = "email"
	OTPChannelSMS   OTPChannel = "sms"
	OTPChannelTOTP  OTPChannel = "totp"
)

// OTPChallenge represents an OTP verification challenge.
type OTPChallenge struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UserUUID  string     `gorm:"type:varchar(36);not null;index" json:"user_uuid"`
	CodeHash  string     `gorm:"type:varchar(255);not null" json:"-"`
	Channel   OTPChannel `gorm:"type:varchar(20);not null" json:"channel"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `gorm:"" json:"used_at,omitempty"`
	IPAddress *string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	Attempts  int        `gorm:"not null;default:0" json:"attempts"`

	// Relationships
	User User `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
}

func (OTPChallenge) TableName() string {
	return "otp_challenges"
}
