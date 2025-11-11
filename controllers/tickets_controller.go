package controllers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/DerylFeyza/freshdesk-automation/models"
	"github.com/DerylFeyza/freshdesk-automation/repository"
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

	if ticketData.GroupID != nil && *ticketData.GroupID == int64(61000171936) && len(ticketData.Attachments) > 0 {
		attachment := ticketData.Attachments[0]
		attachmentText, err := services.AttachmentToText(attachment.AttachmentURL)

		if err != nil {
			fmt.Printf("Error converting attachment to text: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert attachment to text"})
			return
		}

		fmt.Printf("Attachment text: %s\n", string(attachmentText))

		userPemohonRegex := regexp.MustCompile(`(?is)User.*?Pemohon.*?(\d{8}|\d{6})`)
		match := userPemohonRegex.FindStringSubmatch(string(attachmentText))

		if len(match) >= 2 {
			nik := match[1]
			fmt.Println("Found NIK:", nik)

			proactiveData, err := repository.Proactive.CheckProactiveRole(nik)

			fmt.Print("juantcok", proactiveData)

			var roleName string
			if err != nil {
				roleName = "Not Found"
			} else if proactiveData == nil {
				roleName = "Not Found"
			} else {
				data := proactiveData.(*struct {
					EmpID    string
					RoleName string
				})
				roleName = data.RoleName
			}

			ticketStatusLogs, err := repository.TicketStatusLogs.FindByTicketID(response.Ticket.Ticket_id)
			if err != nil {
				fmt.Println("Error fetching ticket status logs:", err)
				return
			}

			if len(ticketStatusLogs) == 0 {
				fmt.Println("No ticket status logs found")
				return
			}

			latestLog := ticketStatusLogs[0]

			ticketLog := &models.ProactiveLogs{
				Ticket_status_update_log_id: latestLog.Ticket_status_update_log_id,
				Emp_id:                      nik,
				Role_name:                   roleName,
			}
			err = repository.Proactive.Create(ticketLog)
			if err != nil {
				fmt.Printf("❌ Failed to create proactive log: %v\n", err)
			} else {
				fmt.Printf("✅ Created proactive log for NIK: %s with role: %s\n", nik, roleName)
			}
		} else {
			fmt.Println("❌ No 6 or 8 digit NIK found after 'User' and 'Pemohon'")
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ticket fetched and saved successfully",
		"ticket":  response,
	})
}
