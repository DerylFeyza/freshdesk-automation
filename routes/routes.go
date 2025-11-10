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
		{
			tickets.POST("", controllers.CreateTicket)                               // Create ticket
			tickets.GET("", controllers.GetAllTickets)                               // Get all tickets (paginated)
			tickets.GET("/:id", controllers.GetTicket)                               // Get ticket by ID
			tickets.GET("/freshdesk/:id", controllers.FindFreshdeskTicket)           // Get by Freshdesk ID
			tickets.POST("/freshdesk/:id", controllers.FindAndInsertFreshdeskTicket) // Get by Freshdesk ID
			tickets.GET("/status/:status", controllers.GetTicketsByStatus)           // Get by status
			tickets.GET("/search", controllers.SearchTickets)                        // Search tickets
			tickets.GET("/stats", controllers.GetTicketStats)                        // Get stats
			tickets.PUT("/:id", controllers.UpdateTicket)                            // Full update
			tickets.PATCH("/:id", controllers.UpdateTicketFields)                    // Partial update
			tickets.DELETE("/:id", controllers.DeleteTicket)                         // Soft delete
		}
	}

	// Webhook endpoint (outside API versioning)
	r.POST("/webhook", controllers.ReceiveTicket)
}
