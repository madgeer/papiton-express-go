package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/order/models"
)

type OrderService struct {
	repo IOrderRepository
}

func NewOrderService(repo IOrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// CreateOrder processes the business logic of creating a shipment order
func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*models.Order, error) {
	custUUID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customerId format")
	}
	serviceUUID, err := uuid.Parse(req.ServiceTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid serviceTypeId format")
	}

	// 1. Dapatkan tarif asli berdasarkan pencarian kota
	tariff, err := s.repo.GetTariff(req.Sender.City, req.Receiver.City, serviceUUID)
	if err != nil {
		return nil, fmt.Errorf("tarif tidak ditemukan untuk rute %s -> %s", req.Sender.City, req.Receiver.City)
	}

	// 2. Kalkulasi biaya
	shippingCost := req.Weight * tariff.PricePerKg
	insuranceFee := 5000.0
	totalPrice := shippingCost + insuranceFee

	// 3. Generate ID & Nomor Resi
	orderID := uuid.New()
	trackingNum := fmt.Sprintf("PPN-%d", time.Now().UnixNano()/1e6)

	order := models.Order{
		ID:             orderID,
		CustomerID:     custUUID,
		ServiceTypeID:  serviceUUID,
		TrackingNumber: trackingNum,
		Weight:         req.Weight,
		Length:         req.Length,
		Width:          req.Width,
		Height:         req.Height,
		ShippingCost:   shippingCost,
		InsuranceFee:   insuranceFee,
		TotalPrice:     totalPrice,
		Status:         models.StatusWaitingPayment,
	}

	addresses := []models.Address{
		{
			ID:         uuid.New(),
			OrderID:    orderID,
			Type:       models.AddressSender,
			Name:       req.Sender.Name,
			Phone:      req.Sender.Phone,
			Address:    req.Sender.Address,
			City:       req.Sender.City,
			Province:   req.Sender.Province,
			PostalCode: req.Sender.PostalCode,
		},
		{
			ID:         uuid.New(),
			OrderID:    orderID,
			Type:       models.AddressReceiver,
			Name:       req.Receiver.Name,
			Phone:      req.Receiver.Phone,
			Address:    req.Receiver.Address,
			City:       req.Receiver.City,
			Province:   req.Receiver.Province,
			PostalCode: req.Receiver.PostalCode,
		},
	}

	// 4. Buat payload outbox event
	eventPayload, err := json.Marshal(map[string]interface{}{
		"orderId":        orderID.String(),
		"customerId":     custUUID.String(),
		"trackingNumber": trackingNum,
		"totalPrice":     totalPrice,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to serialize event payload")
	}

	outboxEvent := models.OrderOutbox{
		ID:            uuid.New(),
		AggregateType: "order",
		AggregateID:   orderID,
		EventType:     "OrderCreated",
		Payload:       eventPayload,
		Status:        "PENDING",
	}

	// 5. Simpan ke database melalui repository
	err = s.repo.SaveOrder(&order, addresses, &outboxEvent)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrder fetches a single order detail
func (s *OrderService) GetOrder(id uuid.UUID) (*models.Order, error) {
	return s.repo.GetOrderByID(id)
}

// ListOrders fetches all orders
func (s *OrderService) ListOrders() ([]models.Order, error) {
	return s.repo.ListOrders()
}

// CompleteOrder updates the status of an order to COMPLETED
func (s *OrderService) CompleteOrder(id uuid.UUID) (*models.Order, error) {
	order, err := s.repo.GetOrderByIDSimple(id)
	if err != nil {
		return nil, fmt.Errorf("order tidak ditemukan")
	}

	err = s.repo.UpdateOrderStatus(id, models.StatusCompleted)
	if err != nil {
		return nil, err
	}

	order.Status = models.StatusCompleted
	return order, nil
}
