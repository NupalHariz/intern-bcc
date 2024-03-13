package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IInformationRepository interface{
	CreateInformation(newInformation *domain.Information) error
}

type InformationRepository struct {
	db *gorm.DB
}

func NewInformationRepository(db *gorm.DB) IInformationRepository {
	return &InformationRepository{db}
}

func (r *InformationRepository) CreateInformation(newInformation *domain.Information) error {
	tx := r.db.Begin()

	err := r.db.Create(newInformation).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
