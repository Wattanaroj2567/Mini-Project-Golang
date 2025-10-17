package models

// RegisterRequest captures incoming data for /api/auth/register.
type RegisterRequest struct {
	Username        string `json:"username"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// LoginRequest captures credentials for /api/auth/login.
type LoginRequest struct {
	Identifier string `json:"identifier"` // accepts username or email
	Password   string `json:"password"`
}

// ForgotPasswordRequest captures the email for reset flow.
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest captures the payload for /api/auth/reset-password.
type ResetPasswordRequest struct {
	Token           string `json:"token"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// AuthResponse is a lightweight DTO returned after successful authentication.
type AuthResponse struct {
	User  *User  `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
}
