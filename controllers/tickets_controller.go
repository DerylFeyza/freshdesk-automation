package controllers

import (
	"net/http"
	"strconv"

	"github.com/DerylFeyza/freshdesk-automation/models"
	"github.com/DerylFeyza/freshdesk-automation/repository"
	"github.com/DerylFeyza/freshdesk-automation/services"
	"github.com/gin-gonic/gin"
)

func FindAndInsertFreshdeskTicket(c *gin.Context) {

	ticketID := c.Param("id")

	response, err := services.UpsertTicketLog(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ticket fetched and saved successfully",
		"ticket":  response,
	})
}

// CreateTicket creates a new ticket
// POST /api/v1/tickets
func CreateTicket(c *gin.Context) {
	var ticket models.Tickets

	if err := c.BindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.Tickets.Create(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// GetTicket retrieves a ticket by ID
// GET /api/v1/tickets/:id
func GetTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := repository.Tickets.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// GetTicketByTicketID retrieves a ticket by Freshdesk ticket ID
// GET /api/v1/tickets/freshdesk/:ticket_id
func GetTicketByTicketID(c *gin.Context) {
	ticketID := c.Param("ticket_id")

	ticket, err := repository.Tickets.FindByTicketID(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// GetAllTickets retrieves all tickets with pagination
// GET /api/v1/tickets?limit=10&offset=0
func GetAllTickets(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tickets, err := repository.Tickets.FindAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tickets": tickets,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetTicketsByStatus retrieves tickets by status
// GET /api/v1/tickets/status/:status?limit=10&offset=0
func GetTicketsByStatus(c *gin.Context) {
	status, err := strconv.Atoi(c.Param("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tickets, err := repository.Tickets.FindByStatus(status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tickets": tickets,
		"status":  status,
		"limit":   limit,
		"offset":  offset,
	})
}

// UpdateTicket updates an existing ticket
// PUT /api/v1/tickets/:id
func UpdateTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// Get existing ticket
	ticket, err := repository.Tickets.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Bind new data
	if err := c.BindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update
	if err := repository.Tickets.Update(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// UpdateTicketFields partially updates ticket fields
// PATCH /api/v1/tickets/:id
func UpdateTicketFields(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.BindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.Tickets.UpdateFields(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully"})
}

// DeleteTicket soft deletes a ticket
// DELETE /api/v1/tickets/:id
func DeleteTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	if err := repository.Tickets.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

// SearchTickets searches tickets by keyword
// GET /api/v1/tickets/search?q=keyword&limit=10&offset=0
func SearchTickets(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search keyword required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tickets, err := repository.Tickets.Search(keyword, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tickets": tickets,
		"keyword": keyword,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetTicketStats returns ticket statistics
// GET /api/v1/tickets/stats
func GetTicketStats(c *gin.Context) {
	total, err := repository.Tickets.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_tickets": total,
	})
}
