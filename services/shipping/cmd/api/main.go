package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/madgeer/papiton-express-go/common/database"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/handler"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/repository"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/service"

	_ "github.com/madgeer/papiton-express-go/services/shipping/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Papiton Express - Shipping Service API
// @version 1.0
// @description REST API for Papiton Express Shipping Service (Milestone 1)
// @host localhost:8084
// @BasePath /

func main() {
	// Initialize database connection from common module
	db := database.InitPostgres()

	// Dependency Injection Setup
	shippingRepo := repository.NewShippingRepository(db)
	shippingService := service.NewShippingService(shippingRepo)
	shippingHandler := handler.NewShippingHandler(shippingService)

	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Shipping API Routes
	r.POST("/shippings/couriers", shippingHandler.CreateCourierHandler)
	r.GET("/shippings/couriers", shippingHandler.ListCouriersHandler)
	r.POST("/shippings/assignments", shippingHandler.AssignCourierHandler)
	r.GET("/shippings/:id", shippingHandler.GetShipmentHandler)
	r.PUT("/shippings/:id/status", shippingHandler.UpdateStatusHandler)

	// Swagger UI Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	r.Run(":" + port)
}
