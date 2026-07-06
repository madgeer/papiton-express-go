package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/madgeer/papiton-express-go/common/database"
	"github.com/madgeer/papiton-express-go/services/order/internal/handler"
	"github.com/madgeer/papiton-express-go/services/order/internal/repository"
	"github.com/madgeer/papiton-express-go/services/order/internal/seed"
	"github.com/madgeer/papiton-express-go/services/order/internal/service"

	_ "github.com/madgeer/papiton-express-go/services/order/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Papiton Express - Order Service API
// @version 1.0
// @description REST API for Papiton Express Order Service (Milestone 1)
// @host localhost:8081
// @BasePath /

func main() {
	// Initialize PostgreSQL from common module
	db := database.InitPostgres()

	// Run seeder
	seed.SeedDatabase(db)

	// Dependency Injection Setup
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Order REST API Routes
	r.POST("/orders", orderHandler.CreateOrderHandler)
	r.GET("/orders/:id", orderHandler.GetOrderHandler)
	r.GET("/orders", orderHandler.ListOrdersHandler)
	r.POST("/orders/:id/complete", orderHandler.CompleteOrderHandler)

	// Swagger UI Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
