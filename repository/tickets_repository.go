package repository

import (
	"errors"

	"github.com/DerylFeyza/freshdesk-automation/database"
	"github.com/DerylFeyza/freshdesk-automation/models"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

var Tickets *TicketRepository

func InitRepository() {
	Tickets = NewTicketRepository()
	TicketStatusLogs = NewTicketStatusLogRepository()
	Proactive = NewProactiveRepository()
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{
		db: database.DB,
	}
}

func (r *TicketRepository) Create(ticket *models.Tickets) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) FindByTicketID(ticketID string) (*models.Tickets, error) {
	var ticket models.Tickets
	err := r.db.Where("ticket_id = ?", ticketID).First(&ticket).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) FindByFreshdeskID(freshdeskID string) (*models.Tickets, error) {
	var ticket models.Tickets
	err := r.db.Where("ticket_freshdesk_id = ?", freshdeskID).First(&ticket).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}
	return &ticket, nil
}

// FindAll retrieves all tickets with optional pagination
func (r *TicketRepository) FindAll(limit, offset int) ([]models.Tickets, error) {
	var tickets []models.Tickets
	query := r.db.Model(&models.Tickets{})

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

// FindByStatus retrieves tickets by status with optional pagination
func (r *TicketRepository) FindByStatus(status, limit, offset int) ([]models.Tickets, error) {
	var tickets []models.Tickets
	query := r.db.Where("status = ?", status)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

// Update updates an existing ticket
func (r *TicketRepository) Update(ticket *models.Tickets) error {
	return r.db.Save(ticket).Error
}

// UpdateFields updates specific fields of a ticket
func (r *TicketRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&models.Tickets{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("ticket not found")
	}
	return nil
}

// Delete soft deletes a ticket (using GORM's soft delete)
func (r *TicketRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Tickets{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("ticket not found")
	}
	return nil
}

// HardDelete permanently deletes a ticket from the database
func (r *TicketRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Tickets{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("ticket not found")
	}
	return nil
}

// Count returns the total number of tickets
func (r *TicketRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Tickets{}).Count(&count).Error
	return count, err
}

// CountByStatus returns the number of tickets with a specific status
func (r *TicketRepository) CountByStatus(status int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Tickets{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// ExistsByTicketID checks if a ticket exists by Freshdesk ticket ID
func (r *TicketRepository) ExistsByTicketID(ticketID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Tickets{}).Where("ticket_id = ?", ticketID).Count(&count).Error
	return count > 0, err
}

// Search searches tickets by subject or description
func (r *TicketRepository) Search(keyword string, limit, offset int) ([]models.Tickets, error) {
	var tickets []models.Tickets
	query := r.db.Where("subject ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) CreateOrUpdate(ticket *models.Tickets) error {

	var existing models.Tickets
	err := r.db.Where("ticket_freshdesk_id = ?", ticket.Ticket_freshdesk_id).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(ticket).Error
	} else if err != nil {
		return err
	}

	ticket.Ticket_id = existing.Ticket_id
	ticket.Created_at = existing.Created_at
	return r.db.Save(ticket).Error
}

// BatchCreate creates multiple tickets at once
func (r *TicketRepository) BatchCreate(tickets []models.Tickets) error {
	return r.db.Create(&tickets).Error
}

// DeleteByTicketID deletes a ticket by its Freshdesk ticket ID
func (r *TicketRepository) DeleteByTicketID(ticketID string) error {
	result := r.db.Where("ticket_id = ?", ticketID).Delete(&models.Tickets{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("ticket not found")
	}
	return nil
}
