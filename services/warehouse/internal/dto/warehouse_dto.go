package dto

type CreateWarehouseRequest struct {
	Name string `json:"name" binding:"required"`
	City string `json:"city" binding:"required"`
}

type CreateRouteRequest struct {
	OriginWarehouseID string `json:"originWarehouseId" binding:"required,uuid"`
	DestinationCity   string `json:"destinationCity" binding:"required"`
	NextWarehouseID   string `json:"nextWarehouseId" binding:"omitempty,uuid"`
}

type CreateMovementRequest struct {
	OrderID         string `json:"orderId" binding:"required,uuid"`
	WarehouseID     string `json:"warehouseId" binding:"required,uuid"`
	Type            string `json:"type" binding:"required,oneof=WAREHOUSE_IN WAREHOUSE_OUT"`
	DestinationCity string `json:"destinationCity" binding:"required_if=Type WAREHOUSE_IN"`
	Notes           string `json:"notes"`
}

type MovementResponse struct {
	MovementID      string  `json:"movementId"`
	OrderID         string  `json:"orderId"`
	WarehouseID     string  `json:"warehouseId"`
	Type            string  `json:"type"`
	NextWarehouseID *string `json:"nextWarehouseId,omitempty"`
	Message         string  `json:"message"`
}
