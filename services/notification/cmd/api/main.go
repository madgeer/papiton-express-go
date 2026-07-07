package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/madgeer/papiton-express-go/common/database"
	"github.com/madgeer/papiton-express-go/services/notification/internal/handler"
	"github.com/madgeer/papiton-express-go/services/notification/internal/repository"
	"github.com/madgeer/papiton-express-go/services/notification/internal/service"

	_ "github.com/madgeer/papiton-express-go/services/notification/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Papiton Express - Notification Service API
// @version 1.0
// @description REST API for Papiton Express Notification Service (Milestone 1)
// @host localhost:8087
// @BasePath /

func main() {
	// Initialize database connection from common module
	db := database.InitPostgres()

	// Dependency Injection Setup
	notifRepo := repository.NewNotificationRepository(db)
	notifService := service.NewNotificationService(notifRepo)
	notifHandler := handler.NewNotificationHandler(notifService)

	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Notification API Routes
	r.POST("/notifications", notifHandler.SendNotificationHandler)
	r.GET("/notifications/:id", notifHandler.GetNotificationHandler)
	r.GET("/notifications/order/:orderId", notifHandler.ListByOrderHandler)

	// Swagger UI Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}
	r.Run(":" + port)
}
