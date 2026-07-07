package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/notification/models"
	"gorm.io/gorm"
)

type INotificationRepository interface {
	CreateNotificationWithOutbox(ctx context.Context, notification *models.Notification, outbox *models.NotificationOutbox) error
	GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	ListNotificationsByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.Notification, error)
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotificationWithOutbox(
	ctx context.Context,
	notification *models.Notification,
	outbox *models.NotificationOutbox,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(notification).Error; err != nil {
			return err
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *NotificationRepository) GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var notif models.Notification
	err := r.db.WithContext(ctx).First(&notif, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notif, nil
}

func (r *NotificationRepository) ListNotificationsByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.Notification, error) {
	var list []models.Notification
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
