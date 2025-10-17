package handlers

import (
	"errors"
	"net/http"

	"github.com/gamegear/users-service/internal/models"
	"github.com/gamegear/users-service/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles all authentication-related HTTP requests.
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles the user registration request (POST /api/auth/register).
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	res, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) || errors.Is(err, services.ErrPasswordMismatch) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// RegisterAdmin handles administrator registration (POST /api/admin/register).
func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	res, err := h.authService.RegisterAdmin(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserAlreadyExists), errors.Is(err, services.ErrPasswordMismatch):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin"})
		}
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Login handles the user login request (POST /api/auth/login).
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// LoginAdmin handles administrator login (POST /api/admin/login).
func (h *AuthHandler) LoginAdmin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	res, err := h.authService.LoginAdmin(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials), errors.Is(err, services.ErrUnauthorizedRole):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}
	c.JSON(http.StatusOK, res)
}

// Logout handles the user logout request (POST /api/auth/logout).
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// LogoutAdmin handles admin logout (POST /api/admin/logout).
func (h *AuthHandler) LogoutAdmin(c *gin.Context) {
	h.Logout(c)
}

// ForgotPassword handles the request to initiate a password reset (POST /api/auth/forgot-password).
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body, email is required"})
		return
	}

	if err := h.authService.ForgotPassword(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process forgot password request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a password reset link has been sent."})
}

// ForgotPasswordAdmin initiates password reset for admins (POST /api/admin/forgot-password).
func (h *AuthHandler) ForgotPasswordAdmin(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body, email is required"})
		return
	}

	if err := h.authService.ForgotPasswordAdmin(c.Request.Context(), req); err != nil {
		if errors.Is(err, services.ErrUnauthorizedRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process forgot password request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "If an admin account with that email exists, a password reset link has been sent."})
}

// ResetPassword handles setting a new password using a reset token (POST /api/auth/reset-password).
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrTokenInvalidOrExpired) || errors.Is(err, services.ErrPasswordMismatch) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully."})
}

// ResetPasswordAdmin completes password reset for admins (POST /api/admin/reset-password).
func (h *AuthHandler) ResetPasswordAdmin(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.authService.ResetPasswordAdmin(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTokenInvalidOrExpired),
			errors.Is(err, services.ErrPasswordMismatch),
			errors.Is(err, services.ErrUnauthorizedRole):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully."})
}
