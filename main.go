package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/DerylFeyza/freshdesk-automation/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	pathDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get working directory: %v", err)
	}

	err = godotenv.Load(filepath.Join(pathDir, ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// postgres.Connect()

	router := gin.Default()
	routes.RegisterRoutes(router)
	router.Run(fmt.Sprintf("%s:%s",
		os.Getenv("HOST"),
		os.Getenv("PORT"),
	))

}
