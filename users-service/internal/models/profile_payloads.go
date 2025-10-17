package models

// ProfileResponse represents the data returned by /api/user/profile.
type ProfileResponse struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	ProfileImage string `json:"profile_image"`
}

// UpdateProfileRequest captures editable profile fields.
type UpdateProfileRequest struct {
	Username        string `json:"username"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ProfileImage    string `json:"profile_image"`
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
	Password        string `json:"password"`
}
