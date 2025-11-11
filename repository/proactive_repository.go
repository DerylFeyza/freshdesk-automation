package repository

import (
	"errors"
	"fmt"

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
		Select("a.emp_id AS emp_id, b.role_name AS role_name").
		Joins("JOIN proactive2.role b ON a.role_id = b.role_id").
		Where("a.emp_id = ?", emp_id).
		Take(&result).
		Error

	fmt.Printf("Proactive Role Check Result: %+v\n", result)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *ProactiveRepository) Create(log *models.ProactiveLogs) error {
	return r.db.Create(log).Error
}
