package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReceiveTicket(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	fmt.Printf("Received webhook payload: %s\n", string(body))
	c.JSON(http.StatusOK, gin.H{"message": "Webhook received successfully"})
}
