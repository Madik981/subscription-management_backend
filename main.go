package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"subscription-management_backend/handlers"
	"subscription-management_backend/models"
)

func main() {
	db := connectDB()

	if err := db.AutoMigrate(&models.Plan{}, &models.User{}, &models.Billing{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	router := gin.Default()
	h := handlers.NewHandler(db)
	h.RegisterRoutes(router)

	port := getEnv("PORT", "8080")
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func connectDB() *gorm.DB {
	dsn := getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=subscription_management port=5432 sslmode=disable TimeZone=UTC")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return db
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
