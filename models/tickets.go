package models

import (
	"time"
)

type Tickets struct {
	Ticket_id           string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Ticket_freshdesk_id string `gorm:"uniqueIndex"`
	Subject             string
	Description         *string `gorm:"type:text;default:null" json:"description,omitempty"`
	Group               *string `gorm:"default:null" json:"group,omitempty"`
	Created_at          time.Time
	Updated_at          time.Time
}

type TicketStatusUpdateLogs struct {
	Ticket_status_update_log_id string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Status                      int
	Created_at                  time.Time
	Ticket_id                   string
	Ticket                      Tickets
}

type ProactiveLogs struct {
	Proactive_log_id            string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Ticket_status_update_log_id string
	TicketStatusUpdateLog       TicketStatusUpdateLogs
	Emp_id                      string
	Role_name                   string
}
