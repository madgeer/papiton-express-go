package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/dto"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/repository"
	"github.com/madgeer/papiton-express-go/services/shipping/models"
)

type ShippingService struct {
	repo repository.IShippingRepository
}

func NewShippingService(repo repository.IShippingRepository) *ShippingService {
	return &ShippingService{repo: repo}
}

func (s *ShippingService) CreateCourier(ctx context.Context, req *dto.CreateCourierRequest) (*models.Courier, error) {
	courier := &models.Courier{
		ID:        uuid.New(),
		Name:      req.Name,
		Phone:     req.Phone,
		Status:    models.CourierAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.CreateCourier(ctx, courier)
	if err != nil {
		return nil, err
	}
	return courier, nil
}

func (s *ShippingService) ListCouriers(ctx context.Context) ([]models.Courier, error) {
	return s.repo.ListCouriers(ctx)
}

func (s *ShippingService) AssignCourier(ctx context.Context, req *dto.AssignShipmentRequest) (*models.Shipment, error) {
	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid orderId format")
	}
	courierUUID, err := uuid.Parse(req.CourierID)
	if err != nil {
		return nil, errors.New("invalid courierId format")
	}

	courier, err := s.repo.GetCourierByID(ctx, courierUUID)
	if err != nil {
		return nil, errors.New("courier not found")
	}

	if courier.Status != models.CourierAvailable {
		return nil, errors.New("courier is busy or offline")
	}

	shipment := &models.Shipment{
		ID:        uuid.New(),
		OrderID:   orderUUID,
		CourierID: &courierUUID,
		Status:    models.ShipmentPendingPickup,
		Notes:     "Courier assigned to pickup package",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.CreateShipmentWithCourierUpdate(ctx, shipment)
	if err != nil {
		return nil, err
	}

	return shipment, nil
}

func (s *ShippingService) GetShipment(ctx context.Context, id uuid.UUID) (*models.Shipment, error) {
	return s.repo.GetShipmentByID(ctx, id)
}

func (s *ShippingService) UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, req *dto.UpdateShipmentStatusRequest) (*models.Shipment, error) {
	shipment, err := s.repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, errors.New("shipment not found")
	}

	oldStatus := shipment.Status
	shipment.Status = models.ShipmentStatus(req.Status)
	shipment.Notes = req.Notes
	shipment.UpdatedAt = time.Now()

	updateCourier := false
	var courierStatus models.CourierStatus

	if shipment.Status == models.ShipmentDelivered || shipment.Status == models.ShipmentFailed {
		updateCourier = true
		courierStatus = models.CourierAvailable
	}

	payloadBytes, err := json.Marshal(map[string]interface{}{
		"shipmentId": shipment.ID.String(),
		"orderId":    shipment.OrderID.String(),
		"oldStatus":  string(oldStatus),
		"newStatus":  string(shipment.Status),
		"notes":      shipment.Notes,
	})
	if err != nil {
		return nil, errors.New("failed to serialize event payload")
	}

	outboxEvent := &models.ShippingOutbox{
		ID:            uuid.New(),
		AggregateType: "shipping",
		AggregateID:   shipment.ID,
		EventType:     "ShipmentStatusUpdated",
		Payload:       payloadBytes,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
	}

	var courierUUID uuid.UUID
	if shipment.CourierID != nil {
		courierUUID = *shipment.CourierID
	}

	err = s.repo.UpdateShipmentStatusWithOutbox(ctx, shipment, outboxEvent, updateCourier, courierUUID, courierStatus)
	if err != nil {
		return nil, err
	}

	return shipment, nil
}
