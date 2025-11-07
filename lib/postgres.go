package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/DerylFeyza/freshdesk-automation/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	DB = db

	err = DB.AutoMigrate(
		&models.Tickets{},
	)
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	fmt.Println("Database connected and migrated successfully")
}
