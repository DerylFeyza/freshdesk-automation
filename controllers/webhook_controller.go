package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DerylFeyza/freshdesk-automation/dto"
	"github.com/DerylFeyza/freshdesk-automation/services"
	"github.com/gin-gonic/gin"
)

func ReceiveTicket(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	var webhook dto.FreshdeskWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
		return
	}
	fmt.Printf("Received webhook payload: %s\n", string(body))

	ticketIDStr := strconv.FormatInt(webhook.Webhook.TicketID, 10)
	ticketData, err := services.GetFreshdeskTicketByID(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := services.UpsertTicketLog(ticketData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook received successfully", "ticket": response})
}
