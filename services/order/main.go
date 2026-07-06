package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

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
	// Inisialisasi koneksi database GORM
	InitDB()

	// Jalankan database seeder untuk data awal
	SeedDatabase(DB)

	// Dependency Injection Setup
	orderRepo := NewOrderRepository(DB)
	orderService := NewOrderService(orderRepo)
	orderHandler := NewOrderHandler(orderService)

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