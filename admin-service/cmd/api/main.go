package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gamegear/admin-service/internal/clients"
	"github.com/gamegear/admin-service/internal/handlers"
	"github.com/gamegear/admin-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warn: no .env file found, relying on environment variables")
	}

	authBaseURL := os.Getenv("ADMIN_AUTH_SERVICE_URL")
	if authBaseURL == "" {
		log.Println("warn: ADMIN_AUTH_SERVICE_URL not set")
	}

	shopBaseURL := os.Getenv("SHOP_SERVICE_URL")
	if shopBaseURL == "" {
		log.Println("warn: SHOP_SERVICE_URL not set")
	}

	authClient := clients.NewAuthClient(authBaseURL)
	shopClient := clients.NewShopServiceClient(shopBaseURL)

	authService := services.NewAuthService(authClient)
	productService := services.NewProductService(shopClient)
	orderService := services.NewOrderService(shopClient)

	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	handlers.RegisterRoutes(router, authHandler, productHandler, orderHandler)

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("admin-service ready on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
