package request

// CreateOrganizationRequest represents the input payload for creating a new organization.
type CreateOrganizationRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=255"`
	Description *string `json:"description" binding:"omitempty"`
	Level       int     `json:"level" binding:"required,min=1"`
	ParentID    *uint   `json:"parent_id" binding:"omitempty"`
}

// UpdateOrganizationRequest represents the input payload for updating an existing organization.
type UpdateOrganizationRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=255"`
	Description *string `json:"description" binding:"omitempty"`
	Status      *string `json:"status" binding:"omitempty,oneof=active suspended"`
	SSOEnforced *bool   `json:"sso_enforced" binding:"omitempty"`
	MFAEnforced *bool   `json:"mfa_enforced" binding:"omitempty"`
}

// BulkRequest represents the payload for processing multiple organizations at once (e.g. Bulk Delete).
type BulkRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
