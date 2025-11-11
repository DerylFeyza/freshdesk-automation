package controllers

import (
	"net/http"

	"github.com/DerylFeyza/freshdesk-automation/services"
	"github.com/gin-gonic/gin"
)

func FindAndInsertFreshdeskTicket(c *gin.Context) {

	ticketID := c.Param("id")

	ticketData, err := services.GetFreshdeskTicketByID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := services.UpsertTicketLog(ticketData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ticket fetched and saved successfully",
		"ticket":  response,
	})
}
