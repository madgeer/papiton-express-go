package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/dto"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/service"
)

type ShippingHandler struct {
	service *service.ShippingService
}

func NewShippingHandler(svc *service.ShippingService) *ShippingHandler {
	return &ShippingHandler{service: svc}
}

// CreateCourierHandler handles POST /shippings/couriers
// @Summary Register a new courier
// @Tags shippings
// @Accept json
// @Produce json
// @Param request body dto.CreateCourierRequest true "Create Courier Payload"
// @Success 201 {object} models.Courier
// @Router /shippings/couriers [post]
func (h *ShippingHandler) CreateCourierHandler(c *gin.Context) {
	var req dto.CreateCourierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courier, err := h.service.CreateCourier(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, courier)
}

// ListCouriersHandler handles GET /shippings/couriers
// @Summary List all couriers
// @Tags shippings
// @Produce json
// @Success 200 {array} models.Courier
// @Router /shippings/couriers [get]
func (h *ShippingHandler) ListCouriersHandler(c *gin.Context) {
	couriers, err := h.service.ListCouriers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, couriers)
}

// AssignCourierHandler handles POST /shippings/assignments
// @Summary Assign a courier to an order (Create Shipment)
// @Tags shippings
// @Accept json
// @Produce json
// @Param request body dto.AssignShipmentRequest true "Assignment Payload"
// @Success 201 {object} models.Shipment
// @Router /shippings/assignments [post]
func (h *ShippingHandler) AssignCourierHandler(c *gin.Context) {
	var req dto.AssignShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shipment, err := h.service.AssignCourier(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shipment)
}

// UpdateStatusHandler handles PUT /shippings/:id/status
// @Summary Update shipment status
// @Tags shippings
// @Accept json
// @Produce json
// @Param id path string true "Shipment UUID"
// @Param request body dto.UpdateShipmentStatusRequest true "Status Payload"
// @Success 200 {object} models.Shipment
// @Router /shippings/{id}/status [put]
func (h *ShippingHandler) UpdateStatusHandler(c *gin.Context) {
	id := c.Param("id")
	shipmentUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid"})
		return
	}

	var req dto.UpdateShipmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shipment, err := h.service.UpdateShipmentStatus(c.Request.Context(), shipmentUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shipment)
}

// GetShipmentHandler handles GET /shippings/:id
// @Summary Get shipment details by ID
// @Tags shippings
// @Produce json
// @Param id path string true "Shipment UUID"
// @Success 200 {object} models.Shipment
// @Router /shippings/{id} [get]
func (h *ShippingHandler) GetShipmentHandler(c *gin.Context) {
	id := c.Param("id")
	shipmentUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid"})
		return
	}

	shipment, err := h.service.GetShipment(c.Request.Context(), shipmentUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shipment)
}
