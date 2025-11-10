package models

type Tickets struct {
	Ticket_uuid string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Ticket_id   string `gorm:"uniqueIndex"`
	Subject     string
	Description string
	Status      int
	Created_at  string
	Updated_at  string
}
