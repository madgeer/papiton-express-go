package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/dto"
	"github.com/madgeer/papiton-express-go/services/tracking/internal/service"
	"github.com/madgeer/papiton-express-go/services/tracking/models"
)

type MockTrackingRepository struct {
	CreateTrackingFunc                   func(ctx context.Context, tracking *models.Tracking, initialEvent *models.TrackingEvent) error
	GetTrackingByNumberFunc              func(ctx context.Context, num string) (*models.Tracking, error)
	GetTrackingByOrderIDFunc             func(ctx context.Context, orderID uuid.UUID) (*models.Tracking, error)
	AddTrackingEventWithStatusUpdateFunc func(ctx context.Context, tracking *models.Tracking, event *models.TrackingEvent, outbox *models.TrackingOutbox) error
}

func (m *MockTrackingRepository) CreateTracking(ctx context.Context, tracking *models.Tracking, initialEvent *models.TrackingEvent) error {
	return m.CreateTrackingFunc(ctx, tracking, initialEvent)
}

func (m *MockTrackingRepository) GetTrackingByNumber(ctx context.Context, num string) (*models.Tracking, error) {
	return m.GetTrackingByNumberFunc(ctx, num)
}

func (m *MockTrackingRepository) GetTrackingByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Tracking, error) {
	return m.GetTrackingByOrderIDFunc(ctx, orderID)
}

func (m *MockTrackingRepository) AddTrackingEventWithStatusUpdate(
	ctx context.Context,
	tracking *models.Tracking,
	event *models.TrackingEvent,
	outbox *models.TrackingOutbox,
) error {
	return m.AddTrackingEventWithStatusUpdateFunc(ctx, tracking, event, outbox)
}

func TestCreateTracking_Success(t *testing.T) {
	mockRepo := &MockTrackingRepository{
		GetTrackingByOrderIDFunc: func(ctx context.Context, orderID uuid.UUID) (*models.Tracking, error) {
			return nil, errors.New("record not found") // Doesn't exist yet
		},
		CreateTrackingFunc: func(ctx context.Context, tracking *models.Tracking, initialEvent *models.TrackingEvent) error {
			return nil
		},
	}

	svc := service.NewTrackingService(mockRepo)

	req := &dto.CreateTrackingRequest{
		OrderID:        uuid.New().String(),
		TrackingNumber: "PPN-12345",
		InitialStatus:  "WAITING_PAYMENT",
	}

	tracking, err := svc.CreateTracking(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tracking.TrackingNumber != "PPN-12345" {
		t.Errorf("Expected tracking number PPN-12345, got %s", tracking.TrackingNumber)
	}

	if tracking.CurrentStatus != "WAITING_PAYMENT" {
		t.Errorf("Expected status WAITING_PAYMENT, got %s", tracking.CurrentStatus)
	}
}

func TestAddTrackingEvent_Success(t *testing.T) {
	trackingID := uuid.New()
	orderID := uuid.New()

	mockRepo := &MockTrackingRepository{
		GetTrackingByNumberFunc: func(ctx context.Context, num string) (*models.Tracking, error) {
			return &models.Tracking{
				ID:             trackingID,
				OrderID:        orderID,
				TrackingNumber: "PPN-12345",
				CurrentStatus:  "WAITING_PAYMENT",
			}, nil
		},
		AddTrackingEventWithStatusUpdateFunc: func(ctx context.Context, tracking *models.Tracking, event *models.TrackingEvent, outbox *models.TrackingOutbox) error {
			if tracking.CurrentStatus != "PICKED_UP" {
				t.Errorf("Expected current status to be PICKED_UP, got %s", tracking.CurrentStatus)
			}
			if outbox.EventType != "TrackingEventAdded" {
				t.Errorf("Expected event type TrackingEventAdded, got %s", outbox.EventType)
			}
			return nil
		},
	}

	svc := service.NewTrackingService(mockRepo)

	req := &dto.AddEventRequest{
		Status:      "PICKED_UP",
		Location:    "Jakarta Hub",
		Description: "Package picked up by courier",
	}

	_, err := svc.AddTrackingEvent(context.Background(), "PPN-12345", req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
