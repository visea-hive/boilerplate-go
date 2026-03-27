package model

import "time"

// LoginFailureReason represents the reason a login attempt failed.
type LoginFailureReason string

const (
	LoginFailureBadPassword LoginFailureReason = "bad_password"
	LoginFailureNoAccount   LoginFailureReason = "no_account"
	LoginFailureLocked      LoginFailureReason = "locked"
	LoginFailureMFAFailed   LoginFailureReason = "mfa_failed"
)

// LoginAttempt represents a record of a user's attempt to log in.
type LoginAttempt struct {
	ID            uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt     time.Time           `gorm:"autoCreateTime;not null" json:"created_at"`
	Identifier    string              `gorm:"type:varchar(255);not null" json:"identifier"`
	IPAddress     string              `gorm:"type:varchar(45);not null" json:"ip_address"`
	UserAgent     *string             `gorm:"type:text" json:"user_agent,omitempty"`
	Success       bool                `gorm:"not null" json:"success"`
	FailureReason *LoginFailureReason `gorm:"type:varchar(50)" json:"failure_reason,omitempty"`
	UserUUID      *string             `gorm:"type:varchar(36);index" json:"user_uuid,omitempty"`

	// Relationships
	User *User `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}
