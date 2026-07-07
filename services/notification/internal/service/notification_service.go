package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/notification/internal/dto"
	"github.com/madgeer/papiton-express-go/services/notification/internal/repository"
	"github.com/madgeer/papiton-express-go/services/notification/models"
)

type NotificationService struct {
	repo repository.INotificationRepository
}

func NewNotificationService(repo repository.INotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) SendNotification(ctx context.Context, req *dto.SendNotificationRequest) (*models.Notification, error) {
	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid orderId format")
	}

	notifID := uuid.New()
	notification := &models.Notification{
		ID:             notifID,
		OrderID:        orderUUID,
		RecipientEmail: req.RecipientEmail,
		RecipientPhone: req.RecipientPhone,
		Type:           req.Type,
		Content:        req.Content,
		Status:         "SENT", // Simulating delivery as successful
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	payloadBytes, err := json.Marshal(map[string]interface{}{
		"notificationId": notifID.String(),
		"orderId":        notification.OrderID.String(),
		"recipientEmail": notification.RecipientEmail,
		"recipientPhone": notification.RecipientPhone,
		"type":           notification.Type,
		"content":        notification.Content,
		"status":         notification.Status,
	})
	if err != nil {
		return nil, errors.New("failed to serialize event payload")
	}

	outboxEvent := &models.NotificationOutbox{
		ID:            uuid.New(),
		AggregateType: "notification",
		AggregateID:   notifID,
		EventType:     "NotificationSent",
		Payload:       payloadBytes,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
	}

	err = s.repo.CreateNotificationWithOutbox(ctx, notification, outboxEvent)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *NotificationService) GetNotification(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	return s.repo.GetNotificationByID(ctx, id)
}

func (s *NotificationService) ListNotificationsByOrder(ctx context.Context, orderID uuid.UUID) ([]models.Notification, error) {
	return s.repo.ListNotificationsByOrderID(ctx, orderID)
}
