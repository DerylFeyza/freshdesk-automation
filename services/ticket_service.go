package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DerylFeyza/freshdesk-automation/dto"
	"github.com/DerylFeyza/freshdesk-automation/models"
	"github.com/DerylFeyza/freshdesk-automation/repository"
	"resty.dev/v3"
)

func GetFreshdeskTicketByID(ticketID string) (dto.Ticket, error) {
	var ticket dto.Ticket

	apiKey := os.Getenv("FRESHDESK_API_KEY")
	domain := os.Getenv("FRESHDESK_DOMAIN")

	client := resty.New()

	resp, err := client.R().
		SetBasicAuth(apiKey, "X").
		SetResult(&ticket).
		Get(fmt.Sprintf("%s/api/v2/tickets/%s", domain, ticketID))
	if err != nil {
		return dto.Ticket{}, err
	}

	if resp.IsError() {
		return dto.Ticket{}, fmt.Errorf("freshdesk API returned status %s", resp.Status())
	}

	return ticket, nil
}

func UpsertTicketLog(ticketData dto.Ticket) (*dto.UpsertTicketResult, error) {

	var groupPtr *string
	if ticketData.GroupID != nil {
		s := strconv.FormatInt(*ticketData.GroupID, 10)
		groupPtr = &s
	}

	ticket := &models.Tickets{
		Ticket_freshdesk_id: strconv.FormatInt(ticketData.ID, 10),
		Subject:             ticketData.Subject,
		Description:         &ticketData.Description,
		Group:               groupPtr,
		Created_at:          time.Now(),
		Updated_at:          time.Now(),
	}

	existingTicket, _ := repository.Tickets.FindByFreshdeskID(ticket.Ticket_freshdesk_id)

	if err := repository.Tickets.CreateOrUpdate(ticket); err != nil {
		return nil, fmt.Errorf("failed to upsert ticket into database: %w", err)
	}

	var attachments []models.Attachments
	if existingTicket == nil {
		if len(ticketData.Attachments) > 0 {
			attachments = make([]models.Attachments, 0, len(ticketData.Attachments))
			for _, att := range ticketData.Attachments {
				attachment := models.Attachments{
					Ticket_id:  ticket.Ticket_id,
					Attachment: att.AttachmentURL,
				}
				attachments = append(attachments, attachment)
			}

			if err := repository.Attachments.BatchCreate(attachments); err != nil {
				return nil, fmt.Errorf("failed to create attachments: %w", err)
			}
		}

		ticketStatus := &models.TicketStatusUpdateLogs{
			Ticket_id:  ticket.Ticket_id,
			Status:     ticketData.Status,
			Created_at: time.Now(),
		}

		if err := repository.TicketStatusLogs.Create(ticketStatus); err != nil {
			return nil, fmt.Errorf("failed to create ticket status log: %w", err)
		}
	}

	return &dto.UpsertTicketResult{
		Ticket:      ticket,
		Attachments: attachments,
		Status:      ticketData.Status,
	}, nil

}
