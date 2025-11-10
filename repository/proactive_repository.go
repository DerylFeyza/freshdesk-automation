package repository

import (
	"errors"

	"github.com/DerylFeyza/freshdesk-automation/database"
	"github.com/DerylFeyza/freshdesk-automation/models"
	"gorm.io/gorm"
)

type ProactiveRepository struct {
	db         *gorm.DB
	Internaldb *gorm.DB
}

var Proactive *ProactiveRepository

func NewProactiveRepository() *ProactiveRepository {
	return &ProactiveRepository{
		db:         database.DB,
		Internaldb: database.InternalDB,
	}
}

func (r *ProactiveRepository) CheckProactiveRole(emp_id string) (interface{}, error) {
	var result struct {
		EmpID    string
		RoleName string
	}

	err := r.Internaldb.
		Table("proactive2.users a").
		Select("a.emp_id, b.role_name").
		Joins("JOIN proactive2.role b ON a.role_id = b.role_id").
		Where("a.emp_id = ?", emp_id).
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return result, nil
}

func (r *ProactiveRepository) Create(log *models.ProactiveLogs) error {
	return r.db.Create(log).Error
}
