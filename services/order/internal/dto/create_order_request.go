package dto

type AddressDetail struct {
	Name       string `json:"name" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Address    string `json:"address" binding:"required"`
	City       string `json:"city" binding:"required"`
	Province   string `json:"province" binding:"required"`
	PostalCode string `json:"postalCode" binding:"required"`
}

type CreateOrderRequest struct {
	CustomerID    string        `json:"customerId" binding:"required,uuid"`
	ServiceTypeID string        `json:"serviceTypeId" binding:"required,uuid"`
	Weight        float64       `json:"weight" binding:"required,gt=0"`
	Length        float64       `json:"length" binding:"required,gt=0"`
	Width         float64       `json:"width" binding:"required,gt=0"`
	Height        float64       `json:"height" binding:"required,gt=0"`
	Sender        AddressDetail `json:"sender" binding:"required"`
	Receiver      AddressDetail `json:"receiver" binding:"required"`
}
