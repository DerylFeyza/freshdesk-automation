package routes

import (
	"github.com/DerylFeyza/freshdesk-automation/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	ticket := r.Group("/api/tickets")
	{
		ticket.GET("/:id", controllers.FindTicket)
	}
	r.POST("/webhook", controllers.ReceiveTicket)

}
