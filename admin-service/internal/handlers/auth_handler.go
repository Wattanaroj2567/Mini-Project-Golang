package handlers

import (
	"net/http"
	"strings"

	"github.com/gamegear/admin-service/internal/models"
	"github.com/gamegear/admin-service/internal/services"
	"github.com/gin-gonic/gin"
)

// compile-time references to DTOs to avoid unused import warnings while logic remains TODO.
var (
	_ models.AdminRegisterRequest
	_ models.AdminLoginRequest
	_ models.AdminForgotPasswordRequest
	_ models.AdminResetPasswordRequest
)

// AuthHandler exposes admin authentication endpoints.
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler constructs AuthHandler.
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/admin/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.AdminRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid register payload", "error": err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": resp})
}

// Login handles POST /api/admin/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid login payload", "error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": resp})
}

// Logout handles POST /api/admin/logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.authService.Logout(c.Request.Context(), extractBearerToken(c.GetHeader("Authorization"))); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "logout successful"})
}

// ForgotPassword handles POST /api/admin/forgot-password.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.AdminForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid forgot-password payload", "error": err.Error()})
		return
	}

	if err := h.authService.ForgotPassword(c.Request.Context(), req); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "success", "message": "password reset email sent"})
}

// ResetPassword handles POST /api/admin/reset-password.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.AdminResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid reset-password payload", "error": err.Error()})
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "password reset successful"})
}

func extractBearerToken(header string) string {
	header = strings.TrimSpace(header)
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return header
}

func handleServiceError(c *gin.Context, err error) {
	type statusCoder interface {
		error
		StatusCode() int
	}

	if sc, ok := err.(statusCoder); ok {
		c.JSON(sc.StatusCode(), gin.H{"status": "error", "message": sc.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error"})
}
