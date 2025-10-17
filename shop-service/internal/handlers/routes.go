package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes wires shop-service routes to provided handlers.
func RegisterRoutes(
	router *gin.Engine,
	productHandler *ProductHandler,
	cartHandler *CartHandler,
	orderHandler *OrderHandler,
	memberMiddleware gin.HandlerFunc,
	adminMiddleware gin.HandlerFunc,
) {
	api := router.Group("/api")

	api.GET("/products", productHandler.List)
	api.GET("/products/:id", productHandler.Get)
	api.POST("/products", adminMiddleware, productHandler.Create)
	api.PUT("/products/:id", adminMiddleware, productHandler.Update)
	api.DELETE("/products/:id", adminMiddleware, productHandler.Delete)

	api.GET("/cart", memberMiddleware, cartHandler.GetCart)
	api.POST("/cart/add", memberMiddleware, cartHandler.AddItem)
	api.PUT("/cart/update", memberMiddleware, cartHandler.UpdateItem)
	api.DELETE("/cart/remove", memberMiddleware, cartHandler.RemoveItem)

	api.POST("/orders", memberMiddleware, orderHandler.Create)
	api.GET("/orders/history", memberMiddleware, orderHandler.History)
	api.GET("/orders/:id", memberMiddleware, orderHandler.GetByID)

	api.GET("/orders", adminMiddleware, orderHandler.ListAll)
	api.PUT("/orders/:id/status", adminMiddleware, orderHandler.UpdateStatus)
}
