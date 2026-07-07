package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/dto"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/repository"
	"github.com/madgeer/papiton-express-go/services/tracking/models"
)

type TrackingService struct {
	repo repository.ITrackingRepository
}

func NewTrackingService(repo repository.ITrackingRepository) *TrackingService {
	return &TrackingService{repo: repo}
}

func (s *TrackingService) CreateTracking(ctx context.Context, req *dto.CreateTrackingRequest) (*models.Tracking, error) {
	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid orderId format")
	}

	// Cek jika tracking dengan orderId ini sudah ada
	existing, err := s.repo.GetTrackingByOrderID(ctx, orderUUID)
	if err == nil && existing != nil {
		return nil, errors.New("tracking for this order already exists")
	}

	trackingID := uuid.New()
	tracking := &models.Tracking{
		ID:             trackingID,
		OrderID:        orderUUID,
		TrackingNumber: req.TrackingNumber,
		CurrentStatus:  req.InitialStatus,
		LastUpdated:    time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	initialEvent := &models.TrackingEvent{
		ID:          uuid.New(),
		TrackingID:  trackingID,
		Status:      req.InitialStatus,
		Location:    "Origin",
		Description: "Order registered and tracking initialized",
		CreatedAt:   time.Now(),
	}

	err = s.repo.CreateTracking(ctx, tracking, initialEvent)
	if err != nil {
		return nil, err
	}

	return tracking, nil
}

func (s *TrackingService) GetTracking(ctx context.Context, trackingNum string) (*models.Tracking, error) {
	return s.repo.GetTrackingByNumber(ctx, trackingNum)
}

func (s *TrackingService) AddTrackingEvent(ctx context.Context, trackingNum string, req *dto.AddEventRequest) (*models.Tracking, error) {
	tracking, err := s.repo.GetTrackingByNumber(ctx, trackingNum)
	if err != nil {
		return nil, errors.New("tracking record not found")
	}

	tracking.CurrentStatus = req.Status
	tracking.LastUpdated = time.Now()
	tracking.UpdatedAt = time.Now()

	event := &models.TrackingEvent{
		ID:          uuid.New(),
		TrackingID:  tracking.ID,
		Status:      req.Status,
		Location:    req.Location,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	payloadBytes, err := json.Marshal(map[string]interface{}{
		"trackingId":     tracking.ID.String(),
		"trackingNumber": tracking.TrackingNumber,
		"orderId":        tracking.OrderID.String(),
		"status":         req.Status,
		"location":       req.Location,
		"description":    req.Description,
		"timestamp":      event.CreatedAt.Unix(),
	})
	if err != nil {
		return nil, errors.New("failed to serialize event payload")
	}

	outboxEvent := &models.TrackingOutbox{
		ID:            uuid.New(),
		AggregateType: "tracking",
		AggregateID:   tracking.ID,
		EventType:     "TrackingEventAdded",
		Payload:       payloadBytes,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
	}

	err = s.repo.AddTrackingEventWithStatusUpdate(ctx, tracking, event, outboxEvent)
	if err != nil {
		return nil, err
	}

	// Muat ulang daftar event setelah penambahan event baru
	return s.repo.GetTrackingByNumber(ctx, trackingNum)
}
