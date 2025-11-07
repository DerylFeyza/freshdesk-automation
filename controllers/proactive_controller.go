package controllers

import (
	"net/http"

	ticket "github.com/DerylFeyza/freshdesk-automation/services"
	"github.com/gin-gonic/gin"
)

func FindTicket(c *gin.Context) {
	ticketID := c.Param("id")

	body, err := ticket.GetTicketByID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}
