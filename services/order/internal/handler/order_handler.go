package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/order/internal/dto"
	"github.com/madgeer/papiton-express-go/services/order/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}

// CreateOrderHandler handles POST /orders
// @Summary Create a new shipment order
// @Description Create a new shipment order, save addresses and outbox event
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "Create Order Request Payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [post]
func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Order successfully created",
		"orderId":        order.ID,
		"trackingNumber": order.TrackingNumber,
		"totalPrice":     order.TotalPrice,
		"status":         order.Status,
	})
}

// GetOrderHandler handles GET /orders/:id
// @Summary Get order details by ID
// @Description Retrieve a single order details preloaded with addresses and service type
// @Tags orders
// @Produce json
// @Param id path string true "Order UUID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderHandler(c *gin.Context) {
	id := c.Param("id")
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid"})
		return
	}

	order, err := h.service.GetOrder(c.Request.Context(), orderUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrdersHandler handles GET /orders
// @Summary List all orders
// @Description Retrieve all orders in the system preloaded with addresses
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} map[string]interface{}
// @Router /orders [get]
func (h *OrderHandler) ListOrdersHandler(c *gin.Context) {
	orders, err := h.service.ListOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// CompleteOrderHandler handles POST /orders/:id/complete
// @Summary Mark an order as completed
// @Description Update the status of an order to COMPLETED
// @Tags orders
// @Produce json
// @Param id path string true "Order UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/complete [post]
func (h *OrderHandler) CompleteOrderHandler(c *gin.Context) {
	id := c.Param("id")
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid"})
		return
	}

	order, err := h.service.CompleteOrder(c.Request.Context(), orderUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order sukses diselesaikan",
		"orderId": order.ID,
		"status":  order.Status,
	})
}
