package handlers

import (
	"net/http"
	"strings"

	"github.com/gamegear/users-service/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a gin middleware for JWT authentication.
// It extracts the token from the "Authorization" header, validates it,
// and sets the userID and userRole in the context for downstream handlers.
func AuthMiddleware(tokenService services.TokenService, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if len(allowedRoles) > 0 && !roleAllowed(claims.Role, allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func roleAllowed(role string, allowed []string) bool {
	for _, r := range allowed {
		if r == role {
			return true
		}
	}
	return false
}

// RegisterRoutes mounts all HTTP routes for the service.
// It sets up public routes under /api/auth and protected routes under /api/user.
func RegisterRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	profileHandler *ProfileHandler,
	tokenService services.TokenService,
) {
	api := router.Group("/api")

	// Authentication routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/logout", AuthMiddleware(tokenService), authHandler.Logout)
	}

	// User profile routes (protected by AuthMiddleware)
	user := api.Group("/user")
	user.Use(AuthMiddleware(tokenService))
	{
		user.GET("/profile", profileHandler.GetProfile)
		user.PUT("/profile", profileHandler.UpdateProfile)
	}

	// Admin-specific authentication routes
	admin := api.Group("/admin")
	{
		admin.POST("/register", authHandler.RegisterAdmin)
		admin.POST("/login", authHandler.LoginAdmin)
		admin.POST("/forgot-password", authHandler.ForgotPasswordAdmin)
		admin.POST("/reset-password", authHandler.ResetPasswordAdmin)
		admin.POST("/logout", AuthMiddleware(tokenService, "admin"), authHandler.LogoutAdmin)
	}
}
