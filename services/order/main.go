package main

import ( 	
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/order/models"
	"gorm.io/gorm"
)

func main() {
	// Inisialisasi koneksi database GORM
	InitDB()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.POST("/orders", CreateOrderHandler)

	r.Run(":8081")
}

func CreateOrderHandler(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Parsing UUID input
	custUUID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customerId format"})
		return
	}
	serviceUUID, err := uuid.Parse(req.ServiceTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid serviceTypeId format"})
		return
	}

	// 2. Generate ID Order dan nomor resi (tracking number) unik
	orderID := uuid.New()
	trackingNum := fmt.Sprintf("PPN-%d", time.Now().UnixNano()/1e6)

	// Perhitungan biaya simulasi (sebelum dipetakan ke data tarif asli)
	shippingCost := req.Weight * 10000.0 // flat Rp 10.000 per kg
	insuranceFee := 5000.0
	totalPrice := shippingCost + insuranceFee

	order := models.Order{
		ID:             orderID,
		CustomerID:     custUUID,
		ServiceTypeID:  serviceUUID,
		TrackingNumber: trackingNum,
		Weight:         req.Weight,
		Length:         req.Length,
		Width:          req.Width,
		Height:         req.Height,
		ShippingCost:   shippingCost,
		InsuranceFee:   insuranceFee,
		TotalPrice:     totalPrice,
		Status:         models.StatusWaitingPayment,
	}

	senderAddress := models.Address{
		ID:         uuid.New(),
		OrderID:    orderID,
		Type:       models.AddressSender,
		Name:       req.Sender.Name,
		Phone:      req.Sender.Phone,
		Address:    req.Sender.Address,
		City:       req.Sender.City,
		Province:   req.Sender.Province,
		PostalCode: req.Sender.PostalCode,
	}

	receiverAddress := models.Address{
		ID:         uuid.New(),
		OrderID:    orderID,
		Type:       models.AddressReceiver,
		Name:       req.Receiver.Name,
		Phone:      req.Receiver.Phone,
		Address:    req.Receiver.Address,
		City:       req.Receiver.City,
		Province:   req.Receiver.Province,
		PostalCode: req.Receiver.PostalCode,
	}

	// 3. Serialisasi payload event OrderCreated untuk Kafka
	eventPayload, err := json.Marshal(map[string]interface{}{
		"orderId":        orderID.String(),
		"customerId":    custUUID.String(),
		"trackingNumber": trackingNum,
		"totalPrice":     totalPrice,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize event payload"})
		return
	}

	outboxEvent := models.OrderOutbox{
		ID:            uuid.New(),
		AggregateType: "order",
		AggregateID:   orderID,
		EventType:     "OrderCreated",
		Payload:       eventPayload,
		Status:        "PENDING",
	}

	// 4. Eksekusi database transaction untuk menjamin data tersimpan secara atomik
	err = DB.Transaction(func(tx *gorm.DB) error {
		// Simpan Order
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		// Simpan Alamat Pengirim
		if err := tx.Create(&senderAddress).Error; err != nil {
			return err
		}
		// Simpan Alamat Penerima
		if err := tx.Create(&receiverAddress).Error; err != nil {
			return err
		}
		// Simpan outbox event
		if err := tx.Create(&outboxEvent).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save order: %v", err)})
		return
	}

	// 5. Kembalikan respons sukses
	c.JSON(http.StatusCreated, gin.H{
		"message":        "Order successfully created",
		"orderId":        orderID,
		"trackingNumber": trackingNum,
		"totalPrice":     totalPrice,
		"status":         order.Status,
	})
}