package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/madgeer/papiton-express-go/common/database"
	"github.com/madgeer/papiton-express-go/services/order/internal/handler"
	"github.com/madgeer/papiton-express-go/services/order/internal/repository"
	"github.com/madgeer/papiton-express-go/services/order/internal/seed"
	"github.com/madgeer/papiton-express-go/services/order/internal/service"
	"github.com/madgeer/papiton-express-go/services/order/models"
	"gorm.io/gorm"
)

// SetupTestRouter initializes GORM, cleans up orders, and registers Gin routes
func SetupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Load .env relative to this test file's location (services/order/cmd/api/)
	_ = godotenv.Load("../../.env")

	// Connect to GORM DB via common module
	db := database.InitPostgres()

	// Clean up tables before testing to keep tests isolated & repeatable
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.OrderOutbox{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Address{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Order{})

	// Seed service types & tariffs
	seed.SeedDatabase(db)

	// Dependency Injection Setup
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	r := gin.Default()
	r.POST("/orders", orderHandler.CreateOrderHandler)
	r.GET("/orders/:id", orderHandler.GetOrderHandler)

	return r, db
}

// TestAPI_CreateOrderSuccess tests the end-to-end flow of POST /orders
func TestAPI_CreateOrderSuccess(t *testing.T) {
	r, db := SetupTestRouter()

	// Use the seeded Service Type ID (Regular: 550e8400-e29b-41d4-a716-446655440001)
	serviceTypeID := "550e8400-e29b-41d4-a716-446655440001"

	payload := map[string]interface{}{
		"customerId":    uuid.New().String(),
		"serviceTypeId": serviceTypeID,
		"weight":        2.5,
		"length":        30,
		"width":         20,
		"height":        10,
		"sender": map[string]string{
			"name":       "Budi",
			"phone":      "081234567890",
			"address":    "Jl. Sudirman 10",
			"city":       "Jakarta",
			"province":   "DKI Jakarta",
			"postalCode": "10110",
		},
		"receiver": map[string]string{
			"name":       "Sari",
			"phone":      "081298765432",
			"address":    "Jl. Setiabudi 20",
			"city":       "Bandung",
			"province":   "Jawa Barat",
			"postalCode": "40111",
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["orderId"]; !ok {
		t.Errorf("Response body missing orderId")
	}

	trackingNum, ok := response["trackingNumber"].(string)
	if !ok || trackingNum == "" {
		t.Errorf("Response body missing or invalid trackingNumber")
	}

	// Verify order exists in the database
	var dbOrder models.Order
	err := db.Preload("Addresses").First(&dbOrder, "tracking_number = ?", trackingNum).Error
	if err != nil {
		t.Fatalf("Failed to find order in database: %v", err)
	}

	if dbOrder.Status != models.StatusWaitingPayment {
		t.Errorf("Expected database order status to be WAITING_PAYMENT, got %s", dbOrder.Status)
	}

	if len(dbOrder.Addresses) != 2 {
		t.Errorf("Expected 2 addresses in database, got %d", len(dbOrder.Addresses))
	}
}

// TestAPI_CreateOrderRouteNotFound tests when the city route is not supported (no tariff)
func TestAPI_CreateOrderRouteNotFound(t *testing.T) {
	r, _ := SetupTestRouter()

	payload := map[string]interface{}{
		"customerId":    uuid.New().String(),
		"serviceTypeId": "550e8400-e29b-41d4-a716-446655440001",
		"weight":        2.5,
		"length":        30,
		"width":         20,
		"height":        10,
		"sender": map[string]string{
			"name":    "Budi",
			"phone":   "081234567890",
			"address": "Jl. Sudirman 10",
			"city":    "Jakarta",
		},
		"receiver": map[string]string{
			"name":    "Sari",
			"phone":   "081298765432",
			"address": "Jl. Setiabudi 20",
			"city":    "UnknownCity", // This route does not exist in tariffs
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions - should fail with 500 or 400 depending on handler
	if w.Code == http.StatusCreated {
		t.Fatalf("Expected creation to fail for unregistered route, but got 201 Created")
	}
}
