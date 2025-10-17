package models

import "time"

// PasswordResetToken stores reset tokens tied to a user account.
type PasswordResetToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;size:255;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides the table name used by GORM.
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
