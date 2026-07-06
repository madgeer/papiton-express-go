package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/payment/models"
	"gorm.io/gorm"
)

type IPaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	UpdatePaymentStatusWithOutbox(ctx context.Context, payment *models.Payment, outbox *models.PaymentOutbox) error
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepository) GetPaymentByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) UpdatePaymentStatusWithOutbox(ctx context.Context, payment *models.Payment, outbox *models.PaymentOutbox) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(payment).Error; err != nil {
			return err
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		return nil
	})
}
