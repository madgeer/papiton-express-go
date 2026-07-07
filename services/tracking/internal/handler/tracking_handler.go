package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/dto"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/service"
)

type TrackingHandler struct {
	service *service.TrackingService
}

func NewTrackingHandler(svc *service.TrackingService) *TrackingHandler {
	return &TrackingHandler{service: svc}
}

// CreateTrackingHandler handles POST /trackings
// @Summary Initialize tracking for an order
// @Tags trackings
// @Accept json
// @Produce json
// @Param request body dto.CreateTrackingRequest true "Create Tracking Payload"
// @Success 201 {object} models.Tracking
// @Router /trackings [post]
func (h *TrackingHandler) CreateTrackingHandler(c *gin.Context) {
	var req dto.CreateTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tracking, err := h.service.CreateTracking(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tracking)
}

// GetTrackingHandler handles GET /trackings/search
// @Summary Search tracking by tracking number
// @Tags trackings
// @Produce json
// @Param q query string true "Tracking Number"
// @Success 200 {object} models.Tracking
// @Failure 404 {object} map[string]interface{}
// @Router /trackings/search [get]
func (h *TrackingHandler) GetTrackingHandler(c *gin.Context) {
	num := c.Query("q")
	if num == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' (tracking number) wajib diisi"})
		return
	}

	tracking, err := h.service.GetTracking(c.Request.Context(), num)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tracking tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, tracking)
}

// AddEventHandler handles POST /trackings/:num/events
// @Summary Add a new tracking event (milestone) and update current status
// @Tags trackings
// @Accept json
// @Produce json
// @Param num path string true "Tracking Number"
// @Param request body dto.AddEventRequest true "Event Payload"
// @Success 200 {object} models.Tracking
// @Router /trackings/{num}/events [post]
func (h *TrackingHandler) AddEventHandler(c *gin.Context) {
	num := c.Param("num")

	var req dto.AddEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tracking, err := h.service.AddTrackingEvent(c.Request.Context(), num, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracking)
}
