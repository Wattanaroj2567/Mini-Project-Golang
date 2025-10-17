package models

import "time"

// Order represents a checkout order
type Order struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	CartID          uint      `json:"cart_id" gorm:"not null"`
	ShippingAddress string    `json:"shipping_address" gorm:"type:text"`
	PaymentMethod   string    `json:"payment_method" gorm:"size:50"`
	Status          string    `json:"status" gorm:"size:50;default:'pending'"`
	Notes           string    `json:"notes" gorm:"type:text"`
	Total           float64   `json:"total" gorm:"type:decimal(10,2);not null"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Foreign key relationships
	Cart       Cart        `json:"cart" gorm:"foreignKey:CartID"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	LineTotal float64 `json:"line_total" gorm:"type:decimal(10,2);not null"`

	// Foreign key relationships
	Order   Order   `json:"order" gorm:"foreignKey:OrderID"`
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

// TableName methods
func (Order) TableName() string {
	return "orders"
}

func (OrderItem) TableName() string {
	return "order_items"
}
