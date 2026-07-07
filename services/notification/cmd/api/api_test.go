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
	"github.com/madgeer/papiton-express-go/services/notification/internal/handler"
	"github.com/madgeer/papiton-express-go/services/notification/internal/repository"
	"github.com/madgeer/papiton-express-go/services/notification/internal/service"
	"github.com/madgeer/papiton-express-go/services/notification/models"
	"gorm.io/gorm"
)

// SetupTestRouter initializes GORM, cleans up tables, and registers Gin routes
func SetupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Load .env relative to this test file's location (services/notification/cmd/api/)
	_ = godotenv.Load("../../.env")

	// Connect to GORM DB via common module
	db := database.InitPostgres()

	// Clean up tables before testing to keep tests isolated & repeatable
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.NotificationOutbox{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Notification{})

	// Dependency Injection Setup
	notifRepo := repository.NewNotificationRepository(db)
	notifService := service.NewNotificationService(notifRepo)
	notifHandler := handler.NewNotificationHandler(notifService)

	r := gin.Default()
	r.POST("/notifications", notifHandler.SendNotificationHandler)
	r.GET("/notifications/:id", notifHandler.GetNotificationHandler)
	r.GET("/notifications/order/:orderId", notifHandler.ListByOrderHandler)

	return r, db
}

// TestAPI_NotificationWorkflow tests sending and retrieving simulated notification records
func TestAPI_NotificationWorkflow(t *testing.T) {
	r, db := SetupTestRouter()

	orderID := uuid.New().String()

	// 1. Send Notification
	payload := map[string]string{
		"orderId":        orderID,
		"recipientEmail": "budi@gmail.com",
		"recipientPhone": "08123456789",
		"type":           "EMAIL",
		"content":        "Your order PPN-101 is being processed.",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/notifications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d. Body: %s", w.Code, w.Body.String())
	}

	var notifRes map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &notifRes)
	notifID := notifRes["id"].(string)

	// 2. Retrieve Specific Notification
	reqGet, _ := http.NewRequest(http.MethodGet, "/notifications/"+notifID, nil)
	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d. Body: %s", wGet.Code, wGet.Body.String())
	}

	// 3. List by Order ID
	reqList, _ := http.NewRequest(http.MethodGet, "/notifications/order/"+orderID, nil)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d. Body: %s", wList.Code, wList.Body.String())
	}

	var listRes []models.Notification
	_ = json.Unmarshal(wList.Body.Bytes(), &listRes)

	if len(listRes) != 1 {
		t.Errorf("Expected 1 notification in list, got %d", len(listRes))
	}

	// 4. Verify DB Outbox
	var dbOutbox models.NotificationOutbox
	err := db.First(&dbOutbox, "aggregate_id = ?", notifID).Error
	if err != nil {
		t.Fatalf("Outbox event not found in database: %v", err)
	}

	if dbOutbox.EventType != "NotificationSent" {
		t.Errorf("Expected outbox event NotificationSent, got %s", dbOutbox.EventType)
	}
}
