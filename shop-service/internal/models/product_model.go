package models

import "time"

// Product represents catalogue data exposed to guests/members via shop-service.
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock       int       `json:"stock" gorm:"not null"`
	CategoryID  uint      `json:"category_id" gorm:"not null"`
	ImageURL    string    `json:"image_url" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName ensures GORM uses the expected table.
func (Product) TableName() string {
	return "products"
}
