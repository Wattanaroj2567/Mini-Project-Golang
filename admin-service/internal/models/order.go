package models

// OrderSummary represents an order listed in admin console.
type OrderSummary struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	Total         float64 `json:"total"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
}

// OrderDetail extends summary with extra information when viewing order detail.
type OrderDetail struct {
	OrderSummary
	ShippingAddress string             `json:"shipping_address"`
	Notes           string             `json:"notes"`
	Items           []OrderItemSummary `json:"items"`
}

// OrderItemSummary represents a single order line.
type OrderItemSummary struct {
	OrderItemID uint    `json:"order_item_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// UpdateOrderStatusRequest captures payload sent from admin-service to shop-service.
type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}
