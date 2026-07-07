package dto

type SendNotificationRequest struct {
	OrderID        string `json:"orderId" binding:"required,uuid"`
	RecipientEmail string `json:"recipientEmail" binding:"required,email"`
	RecipientPhone string `json:"recipientPhone" binding:"required"`
	Type           string `json:"type" binding:"required,oneof=EMAIL SMS WHATSAPP"`
	Content        string `json:"content" binding:"required"`
}
