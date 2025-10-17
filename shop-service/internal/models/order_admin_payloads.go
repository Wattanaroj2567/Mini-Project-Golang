package models

import "time"

// AdminOrderSummary is returned when admin lists all orders.
type AdminOrderSummary struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Status          string    `json:"status"`
	Total           float64   `json:"total"`
	PaymentMethod   string    `json:"payment_method"`
	ShippingAddress string    `json:"shipping_address"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UpdateOrderStatusRequest captures payload for PUT /api/orders/:id/status.
type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

var allowedOrderStatuses = map[string]struct{}{
	"pending":    {},
	"processing": {},
	"shipped":    {},
	"delivered":  {},
	"cancelled":  {},
}

// IsValidOrderStatus reports whether the supplied status is permitted.
func IsValidOrderStatus(status string) bool {
	_, ok := allowedOrderStatuses[status]
	return ok
}
