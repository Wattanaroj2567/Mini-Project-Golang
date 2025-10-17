package handlers

import (
	"net/http"

	"github.com/gamegear/admin-service/internal/models"
	"github.com/gamegear/admin-service/internal/services"
	"github.com/gin-gonic/gin"
)

// OrderHandler exposes admin order management endpoints.
type OrderHandler struct {
	orderService services.OrderService
}

// NewOrderHandler constructs OrderHandler.
func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// List handles GET /api/admin/orders.
func (h *OrderHandler) List(c *gin.Context) {
	orders, err := h.orderService.List(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": orders})
}

// UpdateStatus handles PUT /api/admin/orders/:id/status.
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid order id", "error": err.Error()})
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid order status payload", "error": err.Error()})
		return
	}

	if req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "status is required"})
		return
	}

	if err := h.orderService.UpdateStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "order status updated"})
}
