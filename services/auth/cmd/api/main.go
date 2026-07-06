package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/madgeer/papiton-express-go/common/database"
	"github.com/madgeer/papiton-express-go/services/auth/internal/handler"
	"github.com/madgeer/papiton-express-go/services/auth/internal/repository"
	"github.com/madgeer/papiton-express-go/services/auth/internal/service"

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
	// Initialize database connection from common module
	db := database.InitPostgres()

	// Dependency Injection Setup
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

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
