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
	var count int64
	db.Model(&models.ServiceType{}).Count(&count)
	if count > 0 {
		fmt.Println("Database already has service types, skipping seeding...")
		return
	}

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
		if err := db.Create(&st).Error; err != nil {
			log.Fatalf("Failed to seed service type %s: %v", st.Name, err)
		}
	}
	fmt.Println("Successfully seeded service types!")

	// 2. Seed Tariffs
	tariffs := []models.Tariff{
		{
			ID:              uuid.New(),
			OriginCity:      "Jakarta",
			DestinationCity: "Bandung",
			ServiceTypeID:   regularID,
			PricePerKg:      10000.00,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			OriginCity:      "Jakarta",
			DestinationCity: "Bandung",
			ServiceTypeID:   nextdayID,
			PricePerKg:      20000.00,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			OriginCity:      "Jakarta",
			DestinationCity: "Surabaya",
			ServiceTypeID:   regularID,
			PricePerKg:      25000.00,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			OriginCity:      "Jakarta",
			DestinationCity: "Surabaya",
			ServiceTypeID:   nextdayID,
			PricePerKg:      45000.00,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			OriginCity:      "Bandung",
			DestinationCity: "Surabaya",
			ServiceTypeID:   regularID,
			PricePerKg:      22000.00,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			OriginCity:      "Bandung",
			DestinationCity: "Surabaya",
			ServiceTypeID:   nextdayID,
			PricePerKg:      40000.00,
			CreatedAt:       time.Now(),
		},
	}

	for _, t := range tariffs {
		if err := db.Create(&t).Error; err != nil {
			log.Fatalf("Failed to seed tariff %s -> %s: %v", t.OriginCity, t.DestinationCity, err)
		}
	}
	fmt.Println("Successfully seeded tariffs!")
}
