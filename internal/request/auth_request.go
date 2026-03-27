package request

// RegisterRequest represents the payload for user registration.
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required,min=2"`
}

// VerifyEmailRequest represents the payload for email verification.
type VerifyEmailRequest struct {
	Token string `form:"token" binding:"required"`
}
