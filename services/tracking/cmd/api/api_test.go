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
	"github.com/madgeer/papiton-express-go/services/tracking/internal/handler"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/repository"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/service"
	"github.com/madgeer/papiton-express-go/services/tracking/models"
	"gorm.io/gorm"
)

// SetupTestRouter initializes GORM, cleans up tables, and registers Gin routes
func SetupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Load .env relative to this test file's location (services/tracking/cmd/api/)
	_ = godotenv.Load("../../.env")

	// Connect to GORM DB via common module
	db := database.InitPostgres()

	// Clean up tables before testing to keep tests isolated & repeatable
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.TrackingOutbox{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.TrackingEvent{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Tracking{})

	// Dependency Injection Setup
	trackingRepo := repository.NewTrackingRepository(db)
	trackingService := service.NewTrackingService(trackingRepo)
	trackingHandler := handler.NewTrackingHandler(trackingService)

	r := gin.Default()
	r.POST("/trackings", trackingHandler.CreateTrackingHandler)
	r.GET("/trackings/search", trackingHandler.GetTrackingHandler)
	r.POST("/trackings/:num/events", trackingHandler.AddEventHandler)

	return r, db
}

// TestAPI_TrackingWorkflow tests initializing tracking and adding events to the timeline
func TestAPI_TrackingWorkflow(t *testing.T) {
	r, db := SetupTestRouter()

	orderID := uuid.New().String()
	trackingNum := "PPN-TEST-888"

	// 1. Initialize Tracking
	payload := map[string]string{
		"orderId":        orderID,
		"trackingNumber": trackingNum,
		"initialStatus":  "WAITING_PAYMENT",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/trackings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d. Body: %s", w.Code, w.Body.String())
	}

	var trackingRes map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &trackingRes)
	tID := trackingRes["id"].(string)

	// 2. Add Tracking Event
	eventPayload := map[string]string{
		"status":      "PICKED_UP",
		"location":    "Jakarta Central Hub",
		"description": "Package picked up by courier Anto",
	}
	eventBody, _ := json.Marshal(eventPayload)
	reqEvent, _ := http.NewRequest(http.MethodPost, "/trackings/"+trackingNum+"/events", bytes.NewBuffer(eventBody))
	reqEvent.Header.Set("Content-Type", "application/json")

	wEvent := httptest.NewRecorder()
	r.ServeHTTP(wEvent, reqEvent)

	if wEvent.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d. Body: %s", wEvent.Code, wEvent.Body.String())
	}

	// 3. Search Tracking
	reqSearch, _ := http.NewRequest(http.MethodGet, "/trackings/search?q="+trackingNum, nil)
	wSearch := httptest.NewRecorder()
	r.ServeHTTP(wSearch, reqSearch)

	if wSearch.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d. Body: %s", wSearch.Code, wSearch.Body.String())
	}

	var searchRes models.Tracking
	_ = json.Unmarshal(wSearch.Body.Bytes(), &searchRes)

	if searchRes.CurrentStatus != "PICKED_UP" {
		t.Errorf("Expected current status PICKED_UP, got %s", searchRes.CurrentStatus)
	}

	// Should have 2 events in timeline: Initial status + PICKED_UP
	if len(searchRes.Events) != 2 {
		t.Errorf("Expected 2 events in timeline, got %d", len(searchRes.Events))
	}

	// 4. Verify DB Outbox
	var dbOutbox models.TrackingOutbox
	err := db.First(&dbOutbox, "aggregate_id = ?", tID).Error
	if err != nil {
		t.Fatalf("Outbox event not found in database: %v", err)
	}

	if dbOutbox.EventType != "TrackingEventAdded" {
		t.Errorf("Expected outbox event TrackingEventAdded, got %s", dbOutbox.EventType)
	}
}
