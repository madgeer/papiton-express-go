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
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/handler"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/repository"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/service"
	"github.com/madgeer/papiton-express-go/services/warehouse/models"
	"gorm.io/gorm"
)

// SetupTestRouter initializes GORM, cleans up tables, and registers Gin routes
func SetupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Load .env relative to this test file's location (services/warehouse/cmd/api/)
	_ = godotenv.Load("../../.env")

	// Connect to GORM DB via common module
	db := database.InitPostgres()

	// Clean up tables before testing to keep tests isolated & repeatable
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.WarehouseOutbox{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.WarehouseMovement{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.WarehouseRoute{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Warehouse{})

	// Dependency Injection Setup
	whRepo := repository.NewWarehouseRepository(db)
	whService := service.NewWarehouseService(whRepo)
	whHandler := handler.NewWarehouseHandler(whService)

	r := gin.Default()
	r.POST("/warehouses", whHandler.CreateWarehouseHandler)
	r.GET("/warehouses", whHandler.ListWarehousesHandler)
	r.POST("/warehouses/routes", whHandler.CreateRouteHandler)
	r.POST("/warehouses/movements", whHandler.RecordMovementHandler)

	return r, db
}

// TestAPI_WarehouseWorkflow tests registering warehouses, routes, and movements
func TestAPI_WarehouseWorkflow(t *testing.T) {
	r, db := SetupTestRouter()

	// 1. Create Warehouse 1 (Jakarta)
	wh1Payload := map[string]string{
		"name": "Warehouse Jakarta",
		"city": "Jakarta",
	}
	body1, _ := json.Marshal(wh1Payload)
	req1, _ := http.NewRequest(http.MethodPost, "/warehouses", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d. Body: %s", w1.Code, w1.Body.String())
	}

	var wh1Response map[string]interface{}
	_ = json.Unmarshal(w1.Body.Bytes(), &wh1Response)
	wh1ID := wh1Response["id"].(string)

	// Create Warehouse 2 (Surabaya)
	wh2Payload := map[string]string{
		"name": "Warehouse Surabaya",
		"city": "Surabaya",
	}
	body2, _ := json.Marshal(wh2Payload)
	req2, _ := http.NewRequest(http.MethodPost, "/warehouses", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d. Body: %s", w2.Code, w2.Body.String())
	}

	var wh2Response map[string]interface{}
	_ = json.Unmarshal(w2.Body.Bytes(), &wh2Response)
	wh2ID := wh2Response["id"].(string)

	// 2. Register static route: Jakarta -> Bali (Destination) -> Next Transit: Surabaya
	routePayload := map[string]string{
		"originWarehouseId": wh1ID,
		"destinationCity":   "Bali",
		"nextWarehouseId":   wh2ID,
	}
	routeBody, _ := json.Marshal(routePayload)
	reqRoute, _ := http.NewRequest(http.MethodPost, "/warehouses/routes", bytes.NewBuffer(routeBody))
	reqRoute.Header.Set("Content-Type", "application/json")

	wRoute := httptest.NewRecorder()
	r.ServeHTTP(wRoute, reqRoute)

	if wRoute.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d. Body: %s", wRoute.Code, wRoute.Body.String())
	}

	// 3. Scan WAREHOUSE_IN at Jakarta with Destination Bali (should instruct next routing Surabaya!)
	orderID := uuid.New().String()
	movementPayload := map[string]string{
		"orderId":         orderID,
		"warehouseId":     wh1ID,
		"type":            "WAREHOUSE_IN",
		"destinationCity": "Bali",
		"notes":           "Received from sender Budi",
	}
	movementBody, _ := json.Marshal(movementPayload)
	reqMove, _ := http.NewRequest(http.MethodPost, "/warehouses/movements", bytes.NewBuffer(movementBody))
	reqMove.Header.Set("Content-Type", "application/json")

	wMove := httptest.NewRecorder()
	r.ServeHTTP(wMove, reqMove)

	if wMove.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d. Body: %s", wMove.Code, wMove.Body.String())
	}

	var moveResponse map[string]interface{}
	_ = json.Unmarshal(wMove.Body.Bytes(), &moveResponse)

	nextWH, ok := moveResponse["nextWarehouseId"].(string)
	if !ok || nextWH != wh2ID {
		t.Errorf("Expected next warehouse to be Surabaya (%s), got %s", wh2ID, nextWH)
	}

	// 4. Verify DB Outbox
	var dbOutbox models.WarehouseOutbox
	err := db.First(&dbOutbox, "aggregate_id = ?", moveResponse["movementId"].(string)).Error
	if err != nil {
		t.Fatalf("Outbox event not found in database: %v", err)
	}

	if dbOutbox.EventType != "PackageWarehouseIn" {
		t.Errorf("Expected outbox event PackageWarehouseIn, got %s", dbOutbox.EventType)
	}
}
