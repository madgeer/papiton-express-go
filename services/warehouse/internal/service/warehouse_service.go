package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/dto"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/repository"
	"github.com/madgeer/papiton-express-go/services/warehouse/models"
)

type WarehouseService struct {
	repo repository.IWarehouseRepository
}

func NewWarehouseService(repo repository.IWarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) CreateWarehouse(ctx context.Context, req *dto.CreateWarehouseRequest) (*models.Warehouse, error) {
	wh := &models.Warehouse{
		ID:        uuid.New(),
		Name:      req.Name,
		City:      req.City,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.CreateWarehouse(ctx, wh)
	if err != nil {
		return nil, err
	}
	return wh, nil
}

func (s *WarehouseService) ListWarehouses(ctx context.Context) ([]models.Warehouse, error) {
	return s.repo.ListWarehouses(ctx)
}

func (s *WarehouseService) CreateRoute(ctx context.Context, req *dto.CreateRouteRequest) (*models.WarehouseRoute, error) {
	originUUID, err := uuid.Parse(req.OriginWarehouseID)
	if err != nil {
		return nil, errors.New("invalid originWarehouseId format")
	}

	var nextWarehouseID *uuid.UUID
	if req.NextWarehouseID != "" {
		parsedNext, err := uuid.Parse(req.NextWarehouseID)
		if err != nil {
			return nil, errors.New("invalid nextWarehouseId format")
		}
		nextWarehouseID = &parsedNext
	}

	route := &models.WarehouseRoute{
		ID:                uuid.New(),
		OriginWarehouseID: originUUID,
		DestinationCity:   req.DestinationCity,
		NextWarehouseID:   nextWarehouseID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	err = s.repo.CreateRoute(ctx, route)
	if err != nil {
		return nil, err
	}
	return route, nil
}

func (s *WarehouseService) RecordMovement(ctx context.Context, req *dto.CreateMovementRequest) (*dto.MovementResponse, error) {
	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid orderId format")
	}
	whUUID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		return nil, errors.New("invalid warehouseId format")
	}

	// 1. Verify warehouse exists
	wh, err := s.repo.GetWarehouseByID(ctx, whUUID)
	if err != nil {
		return nil, errors.New("warehouse not found")
	}

	movement := &models.WarehouseMovement{
		ID:          uuid.New(),
		OrderID:     orderUUID,
		WarehouseID: whUUID,
		Type:        models.MovementType(req.Type),
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
	}

	var nextWhIDString *string
	message := fmt.Sprintf("Package recorded as %s at warehouse %s (%s)", req.Type, wh.Name, wh.City)

	// 2. If package enters warehouse, check static routing rule
	if movement.Type == models.MovementIn {
		route, err := s.repo.GetNextWarehouseRoute(ctx, whUUID, req.DestinationCity)
		if err == nil && route != nil {
			if route.NextWarehouseID != nil {
				parsedStr := route.NextWarehouseID.String()
				nextWhIDString = &parsedStr
				message = fmt.Sprintf("Package received at %s. Route next to Warehouse ID %s.", wh.Name, parsedStr)
			} else {
				message = fmt.Sprintf("Package arrived at destination warehouse %s. Ready for courier assignment.", wh.Name)
			}
		} else {
			message = fmt.Sprintf("Package received at %s. Direct local courier delivery to %s.", wh.Name, req.DestinationCity)
		}
	}

	// 3. Create Outbox Event
	eventType := "PackageWarehouseOut"
	if movement.Type == models.MovementIn {
		eventType = "PackageWarehouseIn"
	}

	payloadBytes, err := json.Marshal(map[string]interface{}{
		"movementId":      movement.ID.String(),
		"orderId":         movement.OrderID.String(),
		"warehouseId":     movement.WarehouseID.String(),
		"type":            string(movement.Type),
		"nextWarehouseId": nextWhIDString,
		"notes":           movement.Notes,
	})
	if err != nil {
		return nil, errors.New("failed to serialize event payload")
	}

	outboxEvent := &models.WarehouseOutbox{
		ID:            uuid.New(),
		AggregateType: "warehouse",
		AggregateID:   movement.ID,
		EventType:     eventType,
		Payload:       payloadBytes,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
	}

	// 4. Save to database transaksional
	err = s.repo.SaveMovementWithOutbox(ctx, movement, outboxEvent)
	if err != nil {
		return nil, err
	}

	return &dto.MovementResponse{
		MovementID:      movement.ID.String(),
		OrderID:         movement.OrderID.String(),
		WarehouseID:     movement.WarehouseID.String(),
		Type:            string(movement.Type),
		NextWarehouseID: nextWhIDString,
		Message:         message,
	}, nil
}
