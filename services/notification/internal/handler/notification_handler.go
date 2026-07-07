package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/notification/internal/dto"
	"github.com/madgeer/papiton-express-go/services/notification/internal/service"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: svc}
}

// SendNotificationHandler handles POST /notifications
// @Summary Send and record a simulated notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body dto.SendNotificationRequest true "Send Notification Payload"
// @Success 201 {object} models.Notification
// @Router /notifications [post]
func (h *NotificationHandler) SendNotificationHandler(c *gin.Context) {
	var req dto.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notif, err := h.service.SendNotification(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notif)
}

// GetNotificationHandler handles GET /notifications/:id
// @Summary Get notification details by ID
// @Tags notifications
// @Produce json
// @Param id path string true "Notification UUID"
// @Success 200 {object} models.Notification
// @Router /notifications/{id} [get]
func (h *NotificationHandler) GetNotificationHandler(c *gin.Context) {
	id := c.Param("id")
	notifUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid"})
		return
	}

	notif, err := h.service.GetNotification(c.Request.Context(), notifUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data notifikasi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, notif)
}

// ListByOrderHandler handles GET /notifications/order/:orderId
// @Summary Get all notifications sent for a specific order
// @Tags notifications
// @Produce json
// @Param orderId path string true "Order UUID"
// @Success 200 {array} models.Notification
// @Router /notifications/order/{orderId} [get]
func (h *NotificationHandler) ListByOrderHandler(c *gin.Context) {
	orderId := c.Param("orderId")
	orderUUID, err := uuid.Parse(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Order ID tidak valid"})
		return
	}

	list, err := h.service.ListNotificationsByOrder(c.Request.Context(), orderUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}
