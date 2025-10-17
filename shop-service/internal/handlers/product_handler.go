package handlers

import (
	"net/http"
	"strconv"

	"github.com/gamegear/shop-service/internal/models"
	"github.com/gamegear/shop-service/internal/services"
	"github.com/gin-gonic/gin"
)

// ProductHandler exposes catalogue endpoints.
type ProductHandler struct {
	productService services.ProductService
}

// NewProductHandler constructs ProductHandler.
func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// List handles GET /api/products.
func (h *ProductHandler) List(c *gin.Context) {
	var filter models.ProductFilter

	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid query parameters",
		})
		return
	}

	// Set default values
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 12
	}

	// Validate page and limit
	if filter.Page < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "page must be >= 1",
		})
		return
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "limit must be between 1 and 100",
		})
		return
	}

	// Call service
	products, totalItems, err := h.productService.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := int((totalItems + int64(filter.Limit) - 1) / int64(filter.Limit))

	// Return response with pagination
	c.JSON(http.StatusOK, gin.H{
		"items": products,
		"pagination": gin.H{
			"page":        filter.Page,
			"limit":       filter.Limit,
			"total_items": totalItems,
			"total_pages": totalPages,
		},
	})
}

// Get handles GET /api/products/:id.
func (h *ProductHandler) Get(c *gin.Context) {
	// Parse path param
	idStr := c.Param("id")

	// Convert to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid product ID",
		})
		return
	}

	// Call service
	product, err := h.productService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Not Found",
			Message: "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

// Create handles POST /api/products (admin).
func (h *ProductHandler) Create(c *gin.Context) {
	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	product, err := h.productService.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// Update handles PUT /api/products/:id (admin).
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid product ID",
		})
		return
	}

	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	product, err := h.productService.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// Delete handles DELETE /api/products/:id (admin).
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid product ID",
		})
		return
	}

	if err := h.productService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
