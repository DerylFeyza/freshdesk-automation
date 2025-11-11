package services

import (
	"fmt"
	"os"
	"regexp"
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

	if existingTicket == nil {

		ticketStatus := &models.TicketStatusUpdateLogs{
			Ticket_id:  ticket.Ticket_id,
			Status:     ticketData.Status,
			Created_at: time.Now(),
		}

		if err := repository.TicketStatusLogs.Create(ticketStatus); err != nil {
			return nil, fmt.Errorf("failed to create ticket status log: %w", err)
		}
	}

	ticketStatusLogs, err := repository.TicketStatusLogs.FindByTicketID(ticket.Ticket_id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ticket status logs: %w", err)
	}

	if len(ticketStatusLogs) == 0 {
		return nil, fmt.Errorf("no ticket status logs found")
	}

	latestLog := ticketStatusLogs[0]

	if latestLog.Status != ticketData.Status {
		ticketStatus := &models.TicketStatusUpdateLogs{
			Ticket_id:  ticket.Ticket_id,
			Status:     ticketData.Status,
			Created_at: time.Now(),
		}

		if err := repository.TicketStatusLogs.Create(ticketStatus); err != nil {
			return nil, fmt.Errorf("failed to create ticket status log: %w", err)
		}
	}

	if ticketData.GroupID != nil && *ticketData.GroupID == int64(61000171936) && len(ticketData.Attachments) > 0 {
		attachment := ticketData.Attachments[0]
		attachmentText, err := AttachmentToText(attachment.AttachmentURL)

		if err != nil {

			return nil, fmt.Errorf("error converting attachment to text: %v", err)

		}

		userPemohonRegex := regexp.MustCompile(`(?is)User.*?Pemohon.*?(\d{8}|\d{6})`)
		match := userPemohonRegex.FindStringSubmatch(string(attachmentText))

		if len(match) >= 2 {
			nik := match[1]
			proactiveData, err := repository.Proactive.CheckProactiveRole(nik)
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

			ticketStatusLogs, err := repository.TicketStatusLogs.FindByTicketID(ticket.Ticket_id)
			if err != nil {
				return nil, fmt.Errorf("error fetching ticket status logs: %v", err)
			}

			if len(ticketStatusLogs) == 0 {
				return nil, fmt.Errorf("no ticket status logs found")
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

	return &dto.UpsertTicketResult{
		Ticket: ticket,
		Status: ticketData.Status,
	}, nil

}
