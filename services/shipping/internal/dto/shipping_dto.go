package dto

type CreateCourierRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type AssignShipmentRequest struct {
	OrderID   string `json:"orderId" binding:"required,uuid"`
	CourierID string `json:"courierId" binding:"required,uuid"`
}

type UpdateShipmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=PICKED_UP ON_TRANSIT OUT_FOR_DELIVERY DELIVERED FAILED"`
	Notes  string `json:"notes"`
}
