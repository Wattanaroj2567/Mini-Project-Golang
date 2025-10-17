package models

// ErrorResponse represents error payloads returned by API.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ProductFilter captures query parameters for browsing products.
type ProductFilter struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Category string `form:"category"`
	Search   string `form:"search"`
}

// ProductRequest represents payloads for creating/updating products via admin-service.
type ProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

// ProductResponse serialises product data towards API consumers.
type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}
