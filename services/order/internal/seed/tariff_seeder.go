package seed

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/order/models"
	"gorm.io/gorm"
)

// SeedDatabase inserts initial data for service types and tariffs if they are empty
func SeedDatabase(db *gorm.DB) {
	// 1. Seed Service Types
	regularID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	nextdayID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

	serviceTypes := []models.ServiceType{
		{
			ID:           regularID,
			Name:         "Regular",
			EstimatedDay: 3,
			CreatedAt:    time.Now(),
		},
		{
			ID:           nextdayID,
			Name:         "NextDay",
			EstimatedDay: 1,
			CreatedAt:    time.Now(),
		},
	}

	for _, st := range serviceTypes {
		var existing models.ServiceType
		err := db.First(&existing, "id = ?", st.ID).Error
		if err != nil {
			if err := db.Create(&st).Error; err != nil {
				log.Fatalf("Failed to seed service type %s: %v", st.Name, err)
			}
		}
	}
	fmt.Println("Service types checked/seeded successfully!")

	// 2. Seed Tariffs (both forward and reverse routes)
	tariffs := []models.Tariff{
		// Forward routes
		{OriginCity: "Jakarta", DestinationCity: "Bandung", ServiceTypeID: regularID, PricePerKg: 10000.00},
		{OriginCity: "Jakarta", DestinationCity: "Bandung", ServiceTypeID: nextdayID, PricePerKg: 20000.00},
		{OriginCity: "Jakarta", DestinationCity: "Surabaya", ServiceTypeID: regularID, PricePerKg: 25000.00},
		{OriginCity: "Jakarta", DestinationCity: "Surabaya", ServiceTypeID: nextdayID, PricePerKg: 45000.00},
		{OriginCity: "Bandung", DestinationCity: "Surabaya", ServiceTypeID: regularID, PricePerKg: 22000.00},
		{OriginCity: "Bandung", DestinationCity: "Surabaya", ServiceTypeID: nextdayID, PricePerKg: 40000.00},
		// Reverse routes
		{OriginCity: "Bandung", DestinationCity: "Jakarta", ServiceTypeID: regularID, PricePerKg: 10000.00},
		{OriginCity: "Bandung", DestinationCity: "Jakarta", ServiceTypeID: nextdayID, PricePerKg: 20000.00},
		{OriginCity: "Surabaya", DestinationCity: "Jakarta", ServiceTypeID: regularID, PricePerKg: 25000.00},
		{OriginCity: "Surabaya", DestinationCity: "Jakarta", ServiceTypeID: nextdayID, PricePerKg: 45000.00},
		{OriginCity: "Surabaya", DestinationCity: "Bandung", ServiceTypeID: regularID, PricePerKg: 22000.00},
		{OriginCity: "Surabaya", DestinationCity: "Bandung", ServiceTypeID: nextdayID, PricePerKg: 40000.00},
	}

	for _, t := range tariffs {
		var existing models.Tariff
		err := db.First(&existing, "origin_city = ? AND destination_city = ? AND service_type_id = ?", t.OriginCity, t.DestinationCity, t.ServiceTypeID).Error
		if err != nil { // Tariff does not exist yet, create it
			t.ID = uuid.New()
			t.CreatedAt = time.Now()
			if err := db.Create(&t).Error; err != nil {
				log.Fatalf("Failed to seed tariff %s -> %s: %v", t.OriginCity, t.DestinationCity, err)
			}
		}
	}
	fmt.Println("Tariffs checked/seeded successfully!")
}
