package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gamegear/users-service/internal/handlers"
	"github.com/gamegear/users-service/internal/models"
	"github.com/gamegear/users-service/internal/repositories"
	"github.com/gamegear/users-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// main is the entry point for the application.
// It initializes the database, wires up all dependencies (repositories, services, handlers),
// sets up the HTTP router, and starts the server.
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warn: no .env file found, relying on environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.PasswordResetToken{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Wire Repositories
	userRepo := repositories.NewUserRepository(db)
	passwordResetRepo := repositories.NewPasswordResetRepository(db)

	// Wire Services
	tokenService, err := services.NewTokenService()
	if err != nil {
		log.Fatalf("failed to create token service: %v", err)
	}
	authService := services.NewAuthService(userRepo, passwordResetRepo, tokenService)
	profileService := services.NewProfileService(userRepo)

	// Wire Handlers
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(profileService)

	// Setup Router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Health check route
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Register all application routes
	handlers.RegisterRoutes(router, authHandler, profileHandler, tokenService)

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
