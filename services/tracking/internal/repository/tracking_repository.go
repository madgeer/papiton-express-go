package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/tracking/models"
	"gorm.io/gorm"
)

type ITrackingRepository interface {
	CreateTracking(ctx context.Context, tracking *models.Tracking, initialEvent *models.TrackingEvent) error
	GetTrackingByNumber(ctx context.Context, num string) (*models.Tracking, error)
	GetTrackingByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Tracking, error)
	AddTrackingEventWithStatusUpdate(ctx context.Context, tracking *models.Tracking, event *models.TrackingEvent, outbox *models.TrackingOutbox) error
}

type TrackingRepository struct {
	db *gorm.DB
}

func NewTrackingRepository(db *gorm.DB) ITrackingRepository {
	return &TrackingRepository{db: db}
}

func (r *TrackingRepository) CreateTracking(ctx context.Context, tracking *models.Tracking, initialEvent *models.TrackingEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tracking).Error; err != nil {
			return err
		}
		if err := tx.Create(initialEvent).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *TrackingRepository) GetTrackingByNumber(ctx context.Context, num string) (*models.Tracking, error) {
	var t models.Tracking
	err := r.db.WithContext(ctx).Preload("Events", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&t, "tracking_number = ?", num).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TrackingRepository) GetTrackingByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Tracking, error) {
	var t models.Tracking
	err := r.db.WithContext(ctx).First(&t, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TrackingRepository) AddTrackingEventWithStatusUpdate(
	ctx context.Context,
	tracking *models.Tracking,
	event *models.TrackingEvent,
	outbox *models.TrackingOutbox,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(tracking).Error; err != nil {
			return err
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		return nil
	})
}
