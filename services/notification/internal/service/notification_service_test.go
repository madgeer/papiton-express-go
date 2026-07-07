package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/notification/internal/dto"
	"github.com/madgeer/papiton-express-go/services/notification/internal/service"
	"github.com/madgeer/papiton-express-go/services/notification/models"
)

type MockNotificationRepository struct {
	CreateNotificationWithOutboxFunc func(ctx context.Context, notification *models.Notification, outbox *models.NotificationOutbox) error
	GetNotificationByIDFunc          func(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	ListNotificationsByOrderIDFunc   func(ctx context.Context, orderID uuid.UUID) ([]models.Notification, error)
}

func (m *MockNotificationRepository) CreateNotificationWithOutbox(ctx context.Context, notification *models.Notification, outbox *models.NotificationOutbox) error {
	return m.CreateNotificationWithOutboxFunc(ctx, notification, outbox)
}

func (m *MockNotificationRepository) GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	return m.GetNotificationByIDFunc(ctx, id)
}

func (m *MockNotificationRepository) ListNotificationsByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.Notification, error) {
	return m.ListNotificationsByOrderIDFunc(ctx, orderID)
}

func TestSendNotification_Success(t *testing.T) {
	mockRepo := &MockNotificationRepository{
		CreateNotificationWithOutboxFunc: func(ctx context.Context, n *models.Notification, o *models.NotificationOutbox) error {
			if n.Status != "SENT" {
				t.Errorf("Expected status SENT, got %s", n.Status)
			}
			if o.EventType != "NotificationSent" {
				t.Errorf("Expected event type NotificationSent, got %s", o.EventType)
			}
			return nil
		},
	}

	svc := service.NewNotificationService(mockRepo)

	req := &dto.SendNotificationRequest{
		OrderID:        uuid.New().String(),
		RecipientEmail: "budi@gmail.com",
		RecipientPhone: "08123456789",
		Type:           "EMAIL",
		Content:        "Your order PPN-101 has been processed.",
	}

	notif, err := svc.SendNotification(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if notif.RecipientEmail != "budi@gmail.com" {
		t.Errorf("Expected recipient email budi@gmail.com, got %s", notif.RecipientEmail)
	}
}
