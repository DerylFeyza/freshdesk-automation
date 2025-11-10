package controllers

import (
	"fmt"
	"net/http"
	"regexp"

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
			fmt.Println("AWAEGAEGJNAEIBGIHUA", nik, "WHATTHEHELY")

			count := regexp.MustCompile(regexp.QuoteMeta(nik)).FindAllStringIndex(string(attachmentText), -1)
			if len(count) >= 2 {
				fmt.Printf("✅ Found matching NIK: %s (appears %d times)\n", nik, len(count))
			} else {
				fmt.Printf("⚠️ NIK %s only appears %d time(s), expected at least 2\n", nik, len(count))
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
