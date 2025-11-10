package database

import (
	"fmt"
	"log"
	"os"

	"github.com/DerylFeyza/freshdesk-automation/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	InternalDB *gorm.DB
)

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	DB = db

	mysqlDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("INTERNAL_DB_USER"),
		os.Getenv("INTERNAL_DB_PASSWORD"),
		os.Getenv("INTERNAL_DB_HOST"),
		os.Getenv("INTERNAL_DB_PORT"),
	)

	mysqlDB, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		log.Printf("⚠️  Failed to connect to MySQL: %v", err)
		log.Println("Continuing without internal database...")
	} else {
		InternalDB = mysqlDB
		fmt.Println("Internal Database connected successfully")
	}

	err = DB.AutoMigrate(&models.Tickets{})
	if err != nil {
		log.Fatal("failed to migrate Tickets:", err)
	}

	err = DB.AutoMigrate(&models.TicketStatusUpdateLogs{})
	if err != nil {
		log.Fatal("failed to migrate TicketStatusUpdateLogs:", err)
	}

	err = DB.AutoMigrate(&models.Attachments{})
	if err != nil {
		log.Fatal("failed to migrate Attachments:", err)
	}

	err = DB.AutoMigrate(&models.ProactiveLogs{})
	if err != nil {
		log.Fatal("failed to migrate ProactiveLogs:", err)
	}

	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	fmt.Println("Database connected and migrated successfully")
}
