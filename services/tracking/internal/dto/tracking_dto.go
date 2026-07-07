package dto

type CreateTrackingRequest struct {
	OrderID        string `json:"orderId" binding:"required,uuid"`
	TrackingNumber string `json:"trackingNumber" binding:"required"`
	InitialStatus  string `json:"initialStatus" binding:"required"`
}

type AddEventRequest struct {
	Status      string `json:"status" binding:"required"`
	Location    string `json:"location"`
	Description string `json:"description" binding:"required"`
}
