package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"subscription-management/billing-service/handlers"
)

func main() {
	db := connectDB()

	router := gin.Default()
	h := handlers.NewHandler(
		db,
		getEnv("JWT_SECRET", "super-secret-key"),
		getEnv("ACCOUNTS_SERVICE_URL", "http://localhost:8080"),
		getEnv("INTERNAL_TOKEN", "internal-token"),
	)
	h.RegisterRoutes(router)

	port := getEnv("BILLING_PORT", getEnv("PORT", "8081"))
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func connectDB() *gorm.DB {
	dsn := getEnv("BILLING_DATABASE_URL", "")
	if dsn == "" {
		dsn = getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=subscription_management_billing port=5432 sslmode=disable TimeZone=UTC")
	}

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
