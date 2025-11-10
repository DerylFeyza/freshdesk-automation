package repository

import (
	"errors"

	"github.com/DerylFeyza/freshdesk-automation/database"
	"github.com/DerylFeyza/freshdesk-automation/models"
	"gorm.io/gorm"
)

type TicketStatusLogRepository struct {
	db *gorm.DB
}

var TicketStatusLogs *TicketStatusLogRepository

func NewTicketStatusLogRepository() *TicketStatusLogRepository {
	return &TicketStatusLogRepository{
		db: database.DB,
	}
}

// Create creates a new ticket status update log
func (r *TicketStatusLogRepository) Create(log *models.TicketStatusUpdateLogs) error {
	return r.db.Create(log).Error
}

// FindByID retrieves a ticket status log by its ID
func (r *TicketStatusLogRepository) FindByID(logID string) (*models.TicketStatusUpdateLogs, error) {
	var log models.TicketStatusUpdateLogs
	err := r.db.Where("ticket_status_update_log_id = ?", logID).First(&log).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket status log not found")
		}
		return nil, err
	}
	return &log, nil
}

// FindByTicketID retrieves all status logs for a specific ticket
func (r *TicketStatusLogRepository) FindByTicketID(ticketID string) ([]models.TicketStatusUpdateLogs, error) {
	var logs []models.TicketStatusUpdateLogs
	err := r.db.Where("ticket_id = ?", ticketID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// FindAll retrieves all ticket status logs with optional pagination
func (r *TicketStatusLogRepository) FindAll(limit, offset int) ([]models.TicketStatusUpdateLogs, error) {
	var logs []models.TicketStatusUpdateLogs
	query := r.db.Model(&models.TicketStatusUpdateLogs{})

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// FindLatestByTicketID retrieves the most recent status log for a specific ticket
func (r *TicketStatusLogRepository) FindLatestByTicketID(ticketID string) (*models.TicketStatusUpdateLogs, error) {
	var log models.TicketStatusUpdateLogs
	err := r.db.Where("ticket_id = ?", ticketID).Order("created_at DESC").First(&log).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no status logs found for ticket")
		}
		return nil, err
	}
	return &log, nil
}

// Update updates an existing ticket status log
func (r *TicketStatusLogRepository) Update(log *models.TicketStatusUpdateLogs) error {
	return r.db.Save(log).Error
}

// UpdateFields updates specific fields of a ticket status log
func (r *TicketStatusLogRepository) UpdateFields(logID string, updates map[string]interface{}) error {
	result := r.db.Model(&models.TicketStatusUpdateLogs{}).Where("ticket_status_update_log_id = ?", logID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("ticket status log not found")
	}
	return nil
}

// Delete deletes a ticket status log
func (r *TicketStatusLogRepository) Delete(logID string) error {
	result := r.db.Where("ticket_status_update_log_id = ?", logID).Delete(&models.TicketStatusUpdateLogs{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("ticket status log not found")
	}
	return nil
}

// DeleteByTicketID deletes all status logs for a specific ticket
func (r *TicketStatusLogRepository) DeleteByTicketID(ticketID string) error {
	result := r.db.Where("ticket_id = ?", ticketID).Delete(&models.TicketStatusUpdateLogs{})
	return result.Error
}

// Count returns the total number of ticket status logs
func (r *TicketStatusLogRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.TicketStatusUpdateLogs{}).Count(&count).Error
	return count, err
}

// CountByTicketID returns the number of status logs for a specific ticket
func (r *TicketStatusLogRepository) CountByTicketID(ticketID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.TicketStatusUpdateLogs{}).Where("ticket_id = ?", ticketID).Count(&count).Error
	return count, err
}

// BatchCreate creates multiple ticket status logs at once
func (r *TicketStatusLogRepository) BatchCreate(logs []models.TicketStatusUpdateLogs) error {
	return r.db.Create(&logs).Error
}

// ExistsByID checks if a ticket status log exists by ID
func (r *TicketStatusLogRepository) ExistsByID(logID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.TicketStatusUpdateLogs{}).Where("ticket_status_update_log_id = ?", logID).Count(&count).Error
	return count > 0, err
}

// FindByDateRange retrieves ticket status logs within a date range
func (r *TicketStatusLogRepository) FindByDateRange(startDate, endDate string, limit, offset int) ([]models.TicketStatusUpdateLogs, error) {
	var logs []models.TicketStatusUpdateLogs
	query := r.db.Where("created_at >= ? AND created_at <= ?", startDate, endDate)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}
