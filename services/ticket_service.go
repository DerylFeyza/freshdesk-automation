package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DerylFeyza/freshdesk-automation/dto"
	"github.com/DerylFeyza/freshdesk-automation/models"
	"github.com/DerylFeyza/freshdesk-automation/repository"

	"resty.dev/v3"
)

func GetFreshdeskTicketByID(ticketID string) ([]byte, error) {
	apiKey := os.Getenv("FRESHDESK_API_KEY")
	domain := os.Getenv("FRESHDESK_DOMAIN")

	client := resty.New()

	response, err := client.R().
		SetBasicAuth(apiKey, "X").
		Get(fmt.Sprintf("%s/api/v2/tickets/%s", domain, ticketID))

	if err != nil {
		return nil, err
	}

	return response.Bytes(), nil
}

func UpsertTicketLog(ticketID string) (*models.Tickets, error) {

	ticketData, err := GetFreshdeskTicketByID(ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ticket from Freshdesk: %w", err)
	}

	var freshdeskTicket dto.Ticket
	if err := json.Unmarshal(ticketData, &freshdeskTicket); err != nil {
		return nil, fmt.Errorf("failed to parse Freshdesk response: %w", err)
	}

	ticket := &models.Tickets{
		Ticket_id:   strconv.FormatInt(freshdeskTicket.ID, 10),
		Subject:     freshdeskTicket.Subject,
		Description: freshdeskTicket.DescriptionText,
		Status:      freshdeskTicket.Status,
		Created_at:  time.Now().Format(time.RFC3339),
		Updated_at:  time.Now().Format(time.RFC3339),
	}

	if err := repository.Tickets.CreateOrUpdate(ticket); err != nil {
		return nil, fmt.Errorf("failed to upsert ticket into database: %w", err)
	}

	return ticket, nil
}
