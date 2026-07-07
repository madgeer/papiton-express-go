package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/madgeer/papiton-express-go/common/database"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/handler"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/repository"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/service"

	_ "github.com/madgeer/papiton-express-go/services/warehouse/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Papiton Express - Warehouse Service API
// @version 1.0
// @description REST API for Papiton Express Warehouse Service (Milestone 1)
// @host localhost:8085
// @BasePath /

func main() {
	// Initialize database connection from common module
	db := database.InitPostgres()

	// Dependency Injection Setup
	whRepo := repository.NewWarehouseRepository(db)
	whService := service.NewWarehouseService(whRepo)
	whHandler := handler.NewWarehouseHandler(whService)

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Warehouse API Routes
	r.POST("/warehouses", whHandler.CreateWarehouseHandler)
	r.GET("/warehouses", whHandler.ListWarehousesHandler)
	r.POST("/warehouses/routes", whHandler.CreateRouteHandler)
	r.POST("/warehouses/movements", whHandler.RecordMovementHandler)

	// Swagger UI Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	r.Run(":" + port)
}
