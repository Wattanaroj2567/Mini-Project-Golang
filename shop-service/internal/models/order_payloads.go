package models

// CreateOrderRequest captures payload for POST /api/orders.
type CreateOrderRequest struct {
	CartID          uint   `json:"cart_id"`
	ShippingAddress string `json:"shipping_address"`
	PaymentMethod   string `json:"payment_method"`
	Notes           string `json:"notes"`
}

// OrderSummary is returned in order history responses.
type OrderSummary struct {
	ID            uint    `json:"id"`
	Status        string  `json:"status"`
	Total         float64 `json:"total"`
	PaymentMethod string  `json:"payment_method"`
}

// OrderDetail extends OrderSummary with line items.
type OrderDetail struct {
	OrderSummary
	ShippingAddress string            `json:"shipping_address"`
	Notes           string            `json:"notes"`
	Items           []OrderItemDetail `json:"items"`
}

// OrderItemDetail provides information per order line.
type OrderItemDetail struct {
	OrderItemID uint            `json:"order_item_id"`
	ProductID   uint            `json:"product_id"`
	Quantity    int             `json:"quantity"`
	Price       float64         `json:"price"`
	Product     *ProductSummary `json:"product,omitempty"`
}
