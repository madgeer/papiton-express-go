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
	"github.com/madgeer/papiton-express-go/services/shipping/internal/handler"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/repository"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/service"
	"github.com/madgeer/papiton-express-go/services/shipping/models"
	"gorm.io/gorm"
)

// SetupTestRouter initializes GORM, cleans up tables, and registers Gin routes
func SetupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Load .env relative to this test file's location (services/shipping/cmd/api/)
	_ = godotenv.Load("../../.env")

	// Connect to GORM DB via common module
	db := database.InitPostgres()

	// Clean up tables before testing to keep tests isolated & repeatable
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.ShippingOutbox{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Shipment{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Courier{})

	// Dependency Injection Setup
	shippingRepo := repository.NewShippingRepository(db)
	shippingService := service.NewShippingService(shippingRepo)
	shippingHandler := handler.NewShippingHandler(shippingService)

	r := gin.Default()
	r.POST("/shippings/couriers", shippingHandler.CreateCourierHandler)
	r.GET("/shippings/couriers", shippingHandler.ListCouriersHandler)
	r.POST("/shippings/assignments", shippingHandler.AssignCourierHandler)
	r.PUT("/shippings/:id/status", shippingHandler.UpdateStatusHandler)
	r.GET("/shippings/:id", shippingHandler.GetShipmentHandler)

	return r, db
}

// TestAPI_ShippingWorkflow tests registering a courier, assigning them to a shipment, and delivering it
func TestAPI_ShippingWorkflow(t *testing.T) {
	r, db := SetupTestRouter()

	// 1. Create Courier
	courierPayload := map[string]string{
		"name":  "Budi Courier",
		"phone": "081122334455",
	}
	body, _ := json.Marshal(courierPayload)
	req, _ := http.NewRequest(http.MethodPost, "/shippings/couriers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d. Body: %s", w.Code, w.Body.String())
	}

	var courierResponse map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &courierResponse)

	courierID, ok := courierResponse["id"].(string)
	if !ok || courierID == "" {
		t.Fatalf("Response missing courier ID")
	}

	// 2. Assign Shipment
	orderID := uuid.New().String()
	assignPayload := map[string]string{
		"orderId":   orderID,
		"courierId": courierID,
	}
	assignBody, _ := json.Marshal(assignPayload)
	reqAssign, _ := http.NewRequest(http.MethodPost, "/shippings/assignments", bytes.NewBuffer(assignBody))
	reqAssign.Header.Set("Content-Type", "application/json")

	wAssign := httptest.NewRecorder()
	r.ServeHTTP(wAssign, reqAssign)

	if wAssign.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d. Body: %s", wAssign.Code, wAssign.Body.String())
	}

	var shipmentResponse map[string]interface{}
	_ = json.Unmarshal(wAssign.Body.Bytes(), &shipmentResponse)

	shipmentID, ok := shipmentResponse["id"].(string)
	if !ok || shipmentID == "" {
		t.Fatalf("Response missing shipment ID")
	}

	// Verify Courier Status is now ON_DELIVERY
	type dbCourier struct {
		Status string
	}
	var courier dbCourier
	err := db.Model(&models.Courier{}).Where("id = ?", courierID).First(&courier).Error
	if err != nil {
		t.Fatalf("Courier not found in DB: %v", err)
	}
	if courier.Status != "ON_DELIVERY" {
		t.Errorf("Expected courier status ON_DELIVERY, got %s", courier.Status)
	}

	// 3. Update Shipment Status to DELIVERED
	updatePayload := map[string]string{
		"status": "DELIVERED",
		"notes":  "Delivered safely",
	}
	updateBody, _ := json.Marshal(updatePayload)
	reqUpdate, _ := http.NewRequest(http.MethodPut, "/shippings/"+shipmentID+"/status", bytes.NewBuffer(updateBody))
	reqUpdate.Header.Set("Content-Type", "application/json")

	wUpdate := httptest.NewRecorder()
	r.ServeHTTP(wUpdate, reqUpdate)

	if wUpdate.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d. Body: %s", wUpdate.Code, wUpdate.Body.String())
	}

	// Verify Courier status is now AVAILABLE again
	err = db.Model(&models.Courier{}).Where("id = ?", courierID).First(&courier).Error
	if err != nil {
		t.Fatalf("Courier not found in DB: %v", err)
	}
	if courier.Status != "AVAILABLE" {
		t.Errorf("Expected courier status AVAILABLE, got %s", courier.Status)
	}

	// Verify Outbox Event exists
	var dbOutbox models.ShippingOutbox
	err = db.First(&dbOutbox, "aggregate_id = ?", shipmentID).Error
	if err != nil {
		t.Fatalf("Outbox event not found in database: %v", err)
	}

	if dbOutbox.EventType != "ShipmentStatusUpdated" {
		t.Errorf("Expected outbox event ShipmentStatusUpdated, got %s", dbOutbox.EventType)
	}
}
