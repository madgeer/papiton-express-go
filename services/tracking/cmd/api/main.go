package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/madgeer/papiton-express-go/common/database"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/handler"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/repository"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/service"

	_ "github.com/madgeer/papiton-express-go/services/tracking/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Papiton Express - Tracking Service API
// @version 1.0
// @description REST API for Papiton Express Tracking Service (Milestone 1)
// @host localhost:8086
// @BasePath /

func main() {
	// Initialize database connection from common module
	db := database.InitPostgres()

	// Dependency Injection Setup
	trackingRepo := repository.NewTrackingRepository(db)
	trackingService := service.NewTrackingService(trackingRepo)
	trackingHandler := handler.NewTrackingHandler(trackingService)

	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Tracking API Routes
	r.POST("/trackings", trackingHandler.CreateTrackingHandler)
	r.GET("/trackings/search", trackingHandler.GetTrackingHandler)
	r.POST("/trackings/:num/events", trackingHandler.AddEventHandler)

	// Swagger UI Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}
	r.Run(":" + port)
}
