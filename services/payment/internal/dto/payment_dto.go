package dto

type CreatePaymentRequest struct {
	OrderID       string  `json:"orderId" binding:"required,uuid"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"paymentMethod" binding:"required"`
}

type WebhookPaymentRequest struct {
	PaymentID     string `json:"paymentId" binding:"required,uuid"`
	TransactionID string `json:"transactionId" binding:"required"`
	Status        string `json:"status" binding:"required,oneof=SUCCESS FAILED"`
}
