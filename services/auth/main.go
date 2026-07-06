package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/madgeer/papiton-express-go/services/auth/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Papiton Express - Auth Service API
// @version 1.0
// @description REST API for Papiton Express Auth Service (Milestone 1)
// @host localhost:8082
// @BasePath /

func main() {
	// Initialize database connection
	InitDB()

	// Dependency Injection Setup
	authRepo := NewAuthRepository(DB)
	authService := NewAuthService(authRepo)
	authHandler := NewAuthHandler(authService)

	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Auth API Routes
	r.POST("/auth/register", authHandler.RegisterHandler)
	r.POST("/auth/login", authHandler.LoginHandler)
	r.POST("/auth/refresh", authHandler.RefreshTokenHandler)

	// Swagger UI Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}
