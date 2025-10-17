package models

import "time"

// Cart represents a user's shopping cart
type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Subtotal  float64   `json:"subtotal" gorm:"type:decimal(10,2);default:0"`
	Discount  float64   `json:"discount" gorm:"type:decimal(10,2);default:0"`
	Shipping  float64   `json:"shipping" gorm:"type:decimal(10,2);default:0"`
	Total     float64   `json:"total" gorm:"type:decimal(10,2);default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Foreign key relationship
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartID"`
}

// CartItem represents an item in a cart
type CartItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	CartID    uint    `json:"cart_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	LineTotal float64 `json:"line_total" gorm:"type:decimal(10,2);not null"`

	// Foreign key relationships
	Cart    Cart    `json:"cart" gorm:"foreignKey:CartID"`
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

// TableName methods
func (Cart) TableName() string {
	return "carts"
}

func (CartItem) TableName() string {
	return "cart_items"
}
