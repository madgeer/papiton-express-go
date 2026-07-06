package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/dto"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/service"
)

type WarehouseHandler struct {
	service *service.WarehouseService
}

func NewWarehouseHandler(svc *service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: svc}
}

// CreateWarehouseHandler handles POST /warehouses
// @Summary Register a new warehouse
// @Tags warehouses
// @Accept json
// @Produce json
// @Param request body dto.CreateWarehouseRequest true "Create Warehouse Payload"
// @Success 201 {object} models.Warehouse
// @Router /warehouses [post]
func (h *WarehouseHandler) CreateWarehouseHandler(c *gin.Context) {
	var req dto.CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wh, err := h.service.CreateWarehouse(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wh)
}

// ListWarehousesHandler handles GET /warehouses
// @Summary List all warehouses
// @Tags warehouses
// @Produce json
// @Success 200 {array} models.Warehouse
// @Router /warehouses [get]
func (h *WarehouseHandler) ListWarehousesHandler(c *gin.Context) {
	list, err := h.service.ListWarehouses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// CreateRouteHandler handles POST /warehouses/routes
// @Summary Register a static route rule
// @Tags warehouses
// @Accept json
// @Produce json
// @Param request body dto.CreateRouteRequest true "Create Route Payload"
// @Success 201 {object} models.WarehouseRoute
// @Router /warehouses/routes [post]
func (h *WarehouseHandler) CreateRouteHandler(c *gin.Context) {
	var req dto.CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.service.CreateRoute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, route)
}

// RecordMovementHandler handles POST /warehouses/movements
// @Summary Record package scan movement (IN/OUT) and return routing instructions
// @Tags warehouses
// @Accept json
// @Produce json
// @Param request body dto.CreateMovementRequest true "Record Movement Payload"
// @Success 201 {object} dto.MovementResponse
// @Router /warehouses/movements [post]
func (h *WarehouseHandler) RecordMovementHandler(c *gin.Context) {
	var req dto.CreateMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.RecordMovement(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
