package controllers

import (
	"net/http"

	ticket "github.com/DerylFeyza/freshdesk-automation/services"
	"github.com/gin-gonic/gin"
)

func FindFreshdeskTicket(c *gin.Context) {

	ticketID := c.Param("id")

	body, err := ticket.GetFreshdeskTicketByID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, body)
}
