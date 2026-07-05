package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global GORM database connection instance
var DB *gorm.DB

// InitDB initializes the connection to PostgreSQL for the order database
func InitDB() {
	dsn := "host=localhost user=postgres password=password dbname=papiton_order port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	fmt.Println("Koneksi database GORM berhasil terhubung ke papiton_order!")
}