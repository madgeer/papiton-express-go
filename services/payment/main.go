package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/madgeer/papiton-express-go/services/payment/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Papiton Express - Payment Service API
// @version 1.0
// @description REST API for Papiton Express Payment Service (Milestone 1)
// @host localhost:8083
// @BasePath /

func main() {
	// Initialize database connection
	InitDB()

	// Dependency Injection Setup
	paymentRepo := NewPaymentRepository(DB)
	paymentService := NewPaymentService(paymentRepo)
	paymentHandler := NewPaymentHandler(paymentService)

	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Payment API Routes
	r.POST("/payments", paymentHandler.CreateInvoiceHandler)
	r.GET("/payments/:id", paymentHandler.GetPaymentHandler)
	r.POST("/payments/webhook", paymentHandler.ProcessWebhookHandler)

	// Swagger UI Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	r.Run(":" + port)
}
