package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes wires admin-service endpoints following README specification.
func RegisterRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	productHandler *ProductHandler,
	orderHandler *OrderHandler,
) {
	api := router.Group("/api")

	admin := api.Group("/admin")

	admin.POST("/register", authHandler.Register)
	admin.POST("/login", authHandler.Login)
	admin.POST("/logout", authHandler.Logout)
	admin.POST("/forgot-password", authHandler.ForgotPassword)
	admin.POST("/reset-password", authHandler.ResetPassword)

	products := admin.Group("/products")
	{
		products.GET("", productHandler.List)
		products.POST("", productHandler.Create)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}

	orders := admin.Group("/orders")
	{
		orders.GET("", orderHandler.List)
		orders.PUT("/:id/status", orderHandler.UpdateStatus)
	}
}
