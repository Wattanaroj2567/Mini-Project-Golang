package models

// ProductSummary represents product data exposed on admin console.
type ProductSummary struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

// ProductUpsertRequest is used when admin creates or updates a product.
type ProductUpsertRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}
