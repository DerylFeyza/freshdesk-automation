package routes

import (
	"github.com/DerylFeyza/freshdesk-automation/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Ticket CRUD routes
		tickets := api.Group("/tickets")
		{ // Get ticket by ID
			tickets.GET("/freshdesk/:id", controllers.FindFreshdeskTicket)           // Get by Freshdesk ID
			tickets.POST("/freshdesk/:id", controllers.FindAndInsertFreshdeskTicket) // Get by Freshdesk ID
			// Soft delete
		}
	}

	// Webhook endpoint (outside API versioning)
	r.POST("/webhook", controllers.ReceiveTicket)
}
