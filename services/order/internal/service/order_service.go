package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/order/internal/dto"
	"github.com/madgeer/papiton-express-go/services/order/internal/repository"
	"github.com/madgeer/papiton-express-go/services/order/models"
)

type OrderService struct {
	repo repository.IOrderRepository
}

func NewOrderService(repo repository.IOrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*models.Order, error) {
	custUUID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customerId format")
	}
	serviceUUID, err := uuid.Parse(req.ServiceTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid serviceTypeId format")
	}

	// 1. Dapatkan tarif asli dari repository
	tariff, err := s.repo.GetTariff(ctx, req.Sender.City, req.Receiver.City, serviceUUID)
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
		"customerId":    custUUID.String(),
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
	err = s.repo.SaveOrder(ctx, &order, addresses, &outboxEvent)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s *OrderService) ListOrders(ctx context.Context) ([]models.Order, error) {
	return s.repo.ListOrders(ctx)
}

func (s *OrderService) CompleteOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	order, err := s.repo.GetOrderByIDSimple(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("order tidak ditemukan")
	}

	err = s.repo.UpdateOrderStatus(ctx, id, models.StatusCompleted)
	if err != nil {
		return nil, err
	}

	order.Status = models.StatusCompleted
	return order, nil
}
