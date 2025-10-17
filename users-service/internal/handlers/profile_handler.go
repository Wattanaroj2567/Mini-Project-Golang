package handlers

import (
	"errors"
	"net/http"

	"github.com/gamegear/users-service/internal/models"
	"github.com/gamegear/users-service/internal/services"
	"github.com/gin-gonic/gin"
)

// ProfileHandler handles all user profile-related HTTP requests.
type ProfileHandler struct {
	profileService services.ProfileService
}

// NewProfileHandler creates a new ProfileHandler.
func NewProfileHandler(profileService services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// GetProfile handles the request to view a user's profile (GET /api/user/profile).
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve profile"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// UpdateProfile handles the request to update a user's profile (PUT /api/user/profile).
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedProfile, err := h.profileService.UpdateProfile(c.Request.Context(), userID.(uint), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) || errors.Is(err, services.ErrOldPasswordIncorrect) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrPasswordMismatch) || errors.Is(err, services.ErrPasswordRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	if updatedProfile == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
		return
	}
	c.JSON(http.StatusOK, updatedProfile)
}
