package response

// TokenResponse represents the JWT access and refresh tokens.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Access token TTL in seconds
}

// RegisterResponse represents the response after successful registration.
type RegisterResponse struct {
	UserUUID string `json:"user_uuid"`
	Message  string `json:"message"`
}
