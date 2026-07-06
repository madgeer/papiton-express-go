package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/payment/internal/dto"
	"github.com/madgeer/papiton-express-go/services/payment/internal/service"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: svc}
}

// CreateInvoiceHandler handles POST /payments
// @Summary Create a new payment invoice
// @Description Create a new billing invoice in PENDING status
// @Tags payments
// @Accept json
// @Produce json
// @Param request body dto.CreatePaymentRequest true "Create Invoice Payload"
// @Success 201 {object} models.Payment
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payments [post]
func (h *PaymentHandler) CreateInvoiceHandler(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.CreateInvoice(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// GetPaymentHandler handles GET /payments/:id
// @Summary Get payment status by ID
// @Description Retrieve a single payment record by its ID
// @Tags payments
// @Produce json
// @Param id path string true "Payment UUID"
// @Success 200 {object} models.Payment
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPaymentHandler(c *gin.Context) {
	id := c.Param("id")
	paymentUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid"})
		return
	}

	payment, err := h.service.GetPayment(c.Request.Context(), paymentUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// ProcessWebhookHandler handles POST /payments/webhook
// @Summary Process payment webhook callback
// @Description Webhook from third-party payment gateway to update payment status and record outbox event
// @Tags payments
// @Accept json
// @Produce json
// @Param request body dto.WebhookPaymentRequest true "Webhook Payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payments/webhook [post]
func (h *PaymentHandler) ProcessWebhookHandler(c *gin.Context) {
	var req dto.WebhookPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.ProcessWebhook(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Payment processed successfully",
		"paymentId": payment.ID,
		"status":    payment.Status,
	})
}
