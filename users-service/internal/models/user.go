package models

import "time"

// User represents persisted account information managed by users-service.
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:150;not null"`
	DisplayName  string    `json:"display_name" gorm:"size:150;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;size:255;not null"`
	ProfileImage string    `json:"profile_image" gorm:"size:500"`
	Role         string    `json:"role" gorm:"size:50;default:'member';not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName overrides the table name used by User to match the README specification.
func (User) TableName() string {
	return "users"
}
