package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gamegear/admin-service/internal/clients"
	"github.com/gamegear/admin-service/internal/models"
	"github.com/gamegear/admin-service/internal/services"
	"github.com/gin-gonic/gin"
)

// compile-time references for DTO usage.
var (
	_ models.ProductUpsertRequest
	_ models.ProductSummary
)

// ProductHandler exposes admin product management endpoints.
type ProductHandler struct {
	productService services.ProductService
}

// NewProductHandler constructs ProductHandler.
func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// List handles GET /api/admin/products.
func (h *ProductHandler) List(c *gin.Context) {
	ctx, ok := requireAuthContext(c)
	if !ok {
		return
	}

	products, err := h.productService.List(ctx)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})
}

// Create handles POST /api/admin/products.
func (h *ProductHandler) Create(c *gin.Context) {
	ctx, ok := requireAuthContext(c)
	if !ok {
		return
	}

	var req models.ProductUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid product payload", "error": err.Error()})
		return
	}

	product, err := h.productService.Create(ctx, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": product})
}

// Update handles PUT /api/admin/products/:id.
func (h *ProductHandler) Update(c *gin.Context) {
	ctx, ok := requireAuthContext(c)
	if !ok {
		return
	}

	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid product id", "error": err.Error()})
		return
	}

	var req models.ProductUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid product payload", "error": err.Error()})
		return
	}

	product, err := h.productService.Update(ctx, uint(id), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": product})
}

// Delete handles DELETE /api/admin/products/:id.
func (h *ProductHandler) Delete(c *gin.Context) {
	ctx, ok := requireAuthContext(c)
	if !ok {
		return
	}

	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid product id", "error": err.Error()})
		return
	}

	if err := h.productService.Delete(ctx, uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "product deleted"})
}

func parseUintParam(raw string) (uint64, error) {
	return strconv.ParseUint(raw, 10, 64)
}

func requireAuthContext(c *gin.Context) (context.Context, bool) {
	header := strings.TrimSpace(c.GetHeader("Authorization"))
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Authorization header is required"})
		return nil, false
	}

	return clients.ContextWithAuthToken(c.Request.Context(), header), true
}
