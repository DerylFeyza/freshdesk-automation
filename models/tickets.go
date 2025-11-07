package models

import "gorm.io/gorm"

type Tickets struct {
	gorm.Model
	Ticket_uuid string `gorm:"type:uuid;default:gen_random_uuid()"`
	Ticket_id   string `gorm:"uniqueIndex"`
	Subject     string
	Description string
	Status      int
	Created_at  string
	Updated_at  string
}
