package repository

import (
	"errors"

	"github.com/DerylFeyza/freshdesk-automation/database"
	"github.com/DerylFeyza/freshdesk-automation/models"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

var Attachments *AttachmentRepository

func NewAttachmentRepository() *AttachmentRepository {
	return &AttachmentRepository{
		db: database.DB,
	}
}

func (r *AttachmentRepository) Create(attachment *models.Attachments) error {
	return r.db.Create(attachment).Error
}

func (r *AttachmentRepository) FindByID(attachmentID string) (*models.Attachments, error) {
	var attachment models.Attachments
	err := r.db.Where("attachment_id = ?", attachmentID).First(&attachment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attachment not found")
		}
		return nil, err
	}
	return &attachment, nil
}

func (r *AttachmentRepository) FindByTicketID(ticketID string) ([]models.Attachments, error) {
	var attachments []models.Attachments
	err := r.db.Where("ticket_id = ?", ticketID).Find(&attachments).Error
	return attachments, err
}

func (r *AttachmentRepository) FindAll(limit, offset int) ([]models.Attachments, error) {
	var attachments []models.Attachments
	query := r.db.Model(&models.Attachments{})

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&attachments).Error
	return attachments, err
}

// Update updates an existing attachment
func (r *AttachmentRepository) Update(attachment *models.Attachments) error {
	return r.db.Save(attachment).Error
}

// Delete deletes an attachment
func (r *AttachmentRepository) Delete(attachmentID string) error {
	result := r.db.Where("attachment_id = ?", attachmentID).Delete(&models.Attachments{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("attachment not found")
	}
	return nil
}

// DeleteByTicketID deletes all attachments for a specific ticket
func (r *AttachmentRepository) DeleteByTicketID(ticketID string) error {
	result := r.db.Where("ticket_id = ?", ticketID).Delete(&models.Attachments{})
	return result.Error
}

// Count returns the total number of attachments
func (r *AttachmentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Attachments{}).Count(&count).Error
	return count, err
}

// CountByTicketID returns the number of attachments for a specific ticket
func (r *AttachmentRepository) CountByTicketID(ticketID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Attachments{}).Where("ticket_id = ?", ticketID).Count(&count).Error
	return count, err
}

// BatchCreate creates multiple attachments at once
func (r *AttachmentRepository) BatchCreate(attachments []models.Attachments) error {
	return r.db.Create(&attachments).Error
}

// ExistsByID checks if an attachment exists by ID
func (r *AttachmentRepository) ExistsByID(attachmentID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Attachments{}).Where("attachment_id = ?", attachmentID).Count(&count).Error
	return count > 0, err
}
