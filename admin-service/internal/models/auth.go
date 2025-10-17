package models

// AdminRegisterRequest captures payload for creating a new admin account.
type AdminRegisterRequest struct {
	Email           string `json:"email"`
	DisplayName     string `json:"display_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// AdminLoginRequest captures login credentials.
type AdminLoginRequest struct {
	Email      string `json:"email"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// AdminForgotPasswordRequest captures email for reset token.
type AdminForgotPasswordRequest struct {
	Email string `json:"email"`
}

// AdminResetPasswordRequest captures data for resetting passwords.
type AdminResetPasswordRequest struct {
	Token           string `json:"token"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// AdminAuthResponse describes successful auth response.
type AdminAuthResponse struct {
	Token string       `json:"token"`
	Admin AdminProfile `json:"admin"`
}

// AdminProfile describes admin user info returned to frontend.
type AdminProfile struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}
