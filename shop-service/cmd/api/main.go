package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gamegear/shop-service/internal/handlers"
	"github.com/gamegear/shop-service/internal/middleware"
	"github.com/gamegear/shop-service/internal/models"
	"github.com/gamegear/shop-service/internal/repositories"
	"github.com/gamegear/shop-service/internal/security"
	"github.com/gamegear/shop-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables for local development convenience.
	if err := godotenv.Load(); err != nil {
		log.Println("warn: no .env file found, relying on environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Establish database connection.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate schemas that correspond to README endpoints.
	if err := db.AutoMigrate(
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Wire repositories.
	productRepo := repositories.NewProductRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	cartItemRepo := repositories.NewCartItemRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Wire services.
	productService := services.NewProductService(productRepo)
	cartService := services.NewCartService(cartRepo, cartItemRepo, productRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo, cartItemRepo)

	// Wire handlers.
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService)

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY is required")
	}

	tokenValidator, err := security.NewValidator(jwtSecret)
	if err != nil {
		log.Fatalf("failed to create token validator: %v", err)
	}

	memberMiddleware := middleware.RequireRoles(tokenValidator)
	adminMiddleware := middleware.RequireRoles(tokenValidator, "admin")

	// Setup Gin router with logging/recovery middleware.
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Health check endpoint.
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Register REST endpoints that align with README specification.
	handlers.RegisterRoutes(router, productHandler, cartHandler, orderHandler, memberMiddleware, adminMiddleware)

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("shop-service ready on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
