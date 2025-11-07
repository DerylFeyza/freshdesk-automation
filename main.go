package main

import (
	"fmt"
	"log"
	"os"

	postgres "github.com/DerylFeyza/freshdesk-automation/lib"
	"github.com/DerylFeyza/freshdesk-automation/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgres.Connect()

	router := gin.Default()
	routes.RegisterRoutes(router)
	router.Run(fmt.Sprintf("%s:%s",
		os.Getenv("HOST"),
		os.Getenv("PORT"),
	))

}
