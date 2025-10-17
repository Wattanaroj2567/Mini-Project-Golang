package models

// AddToCartRequest captures payload for POST /api/cart/add.
type AddToCartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// UpdateCartItemRequest captures payload for PUT /api/cart/update.
type UpdateCartItemRequest struct {
	CartItemID uint `json:"cart_item_id"`
	Quantity   int  `json:"quantity"`
}

// RemoveCartItemRequest captures payload for DELETE /api/cart/remove.
type RemoveCartItemRequest struct {
	CartItemID uint `json:"cart_item_id"`
}

// CartResponse serialises cart details including nested items.
type CartResponse struct {
	ID       uint              `json:"id"`
	UserID   uint              `json:"user_id"`
	Items    []CartItemSummary `json:"items"`
	Subtotal float64           `json:"subtotal"`
	Total    float64           `json:"total"`
}

// CartItemSummary provides a lightweight view of each item in cart.
type CartItemSummary struct {
	CartItemID uint            `json:"cart_item_id"`
	ProductID  uint            `json:"product_id"`
	Quantity   int             `json:"quantity"`
	Product    *ProductSummary `json:"product,omitempty"`
}

// ProductSummary nests essential product information inside cart/order responses.
type ProductSummary struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
}
