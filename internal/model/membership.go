package model

import "time"

// MembershipStatus represents the status of a membership.
type MembershipStatus string

const (
	MembershipStatusActive    MembershipStatus = "active"
	MembershipStatusInvited   MembershipStatus = "invited"
	MembershipStatusSuspended MembershipStatus = "suspended"
)

// Membership represents a user's membership in an organization.
type Membership struct {
	BaseModel

	UserUUID string           `gorm:"type:varchar(36);not null;index" json:"user_uuid"`
	OrgID    uint             `gorm:"not null;index" json:"org_id"`
	Status   MembershipStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// Relationships
	User         User         `gorm:"foreignKey:UserUUID;references:UUID" json:"user,omitempty"`
	Organization Organization `gorm:"foreignKey:OrgID;references:ID" json:"organization,omitempty"`
}

// TableName specifies the table name for the Membership model.
func (Membership) TableName() string {
	return "memberships"
}

// MembershipDTO represents the API response for a membership.
type MembershipDTO struct {
	ID           uint            `json:"id"`
	UserUUID     string          `json:"user_uuid"`
	OrgID        uint            `json:"org_id"`
	Status       string          `json:"status"`
	User         UserDTO         `json:"user"`
	Organization OrganizationDTO `json:"organization"`
	CreatedAt    time.Time       `json:"created_at"`
}

func MapToMembershipDTO(m *Membership) *MembershipDTO {
	if m == nil {
		return nil
	}

	var userDTO UserDTO
	var organizationDTO OrganizationDTO

	if m.User.ID != 0 {
		userDTO = *MapToUserDTO(&m.User)
	}

	if m.Organization.ID != 0 {
		organizationDTO = MapToOrganizationDTO(&m.Organization)
	}

	return &MembershipDTO{
		ID:           m.ID,
		UserUUID:     m.UserUUID,
		OrgID:        m.OrgID,
		Status:       string(m.Status),
		User:         userDTO,
		Organization: organizationDTO,
		CreatedAt:    m.CreatedAt,
	}
}

func MapToMembershipDTOs(memberships []Membership) []MembershipDTO {
	var dtos []MembershipDTO
	for _, m := range memberships {
		dto := MapToMembershipDTO(&m)
		dtos = append(dtos, *dto)
	}

	return dtos
}
