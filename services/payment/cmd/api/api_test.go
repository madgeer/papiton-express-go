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
	"github.com/madgeer/papiton-express-go/services/payment/internal/handler"
	"github.com/madgeer/papiton-express-go/services/payment/internal/repository"
	"github.com/madgeer/papiton-express-go/services/payment/internal/service"
	"github.com/madgeer/papiton-express-go/services/payment/models"
	"gorm.io/gorm"
)

// SetupTestRouter initializes GORM, cleans up payments, and registers Gin routes
func SetupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Load .env relative to this test file's location (services/payment/cmd/api/)
	_ = godotenv.Load("../../.env")

	// Connect to GORM DB via common module
	db := database.InitPostgres()

	// Clean up tables before testing to keep tests isolated & repeatable
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.PaymentOutbox{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Payment{})

	// Dependency Injection Setup
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	r := gin.Default()
	r.POST("/payments", paymentHandler.CreateInvoiceHandler)
	r.GET("/payments/:id", paymentHandler.GetPaymentHandler)
	r.POST("/payments/webhook", paymentHandler.ProcessWebhookHandler)

	return r, db
}

// TestAPI_CreateAndVerifyPayment tests end-to-end invoice creation and webhook verification
func TestAPI_CreateAndVerifyPayment(t *testing.T) {
	r, db := SetupTestRouter()

	orderID := uuid.New().String()

	// 1. Create Invoice
	payload := map[string]interface{}{
		"orderId":       orderID,
		"amount":        35000.0,
		"paymentMethod": "SHOPEEPAY",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d. Body: %s", w.Code, w.Body.String())
	}

	var paymentResponse map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &paymentResponse)

	paymentID, ok := paymentResponse["id"].(string)
	if !ok || paymentID == "" {
		t.Fatalf("Response missing payment ID")
	}

	// 2. Webhook Callback (Simulate Successful Payment)
	webhookPayload := map[string]interface{}{
		"paymentId":     paymentID,
		"transactionId": "TX-PAY-88899",
		"status":        "SUCCESS",
	}

	webhookBody, _ := json.Marshal(webhookPayload)
	webhookReq, _ := http.NewRequest(http.MethodPost, "/payments/webhook", bytes.NewBuffer(webhookBody))
	webhookReq.Header.Set("Content-Type", "application/json")

	wWebhook := httptest.NewRecorder()
	r.ServeHTTP(wWebhook, webhookReq)

	if wWebhook.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d. Body: %s", wWebhook.Code, wWebhook.Body.String())
	}

	// 3. Verify Database State
	var dbPayment models.Payment
	err := db.First(&dbPayment, "id = ?", paymentID).Error
	if err != nil {
		t.Fatalf("Payment record not found in database: %v", err)
	}

	if dbPayment.Status != models.StatusSuccess {
		t.Errorf("Expected status SUCCESS, got %s", dbPayment.Status)
	}

	if *dbPayment.TransactionID != "TX-PAY-88899" {
		t.Errorf("Expected TransactionID TX-PAY-88899, got %s", *dbPayment.TransactionID)
	}

	// Verify Outbox Event exists
	var dbOutbox models.PaymentOutbox
	err = db.First(&dbOutbox, "aggregate_id = ?", paymentID).Error
	if err != nil {
		t.Fatalf("Outbox event not found in database: %v", err)
	}

	if dbOutbox.EventType != "PaymentVerified" {
		t.Errorf("Expected outbox event PaymentVerified, got %s", dbOutbox.EventType)
	}
}
