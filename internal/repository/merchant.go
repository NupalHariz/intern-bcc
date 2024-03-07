package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IMerchantRepository interface {
	CreateMerchant(newMerchant *domain.Merchants) error
}

type MerchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) IMerchantRepository {
	return &MerchantRepository{db}
}

func (r *MerchantRepository) CreateMerchant(newMerchant *domain.Merchants) error {
	tx := r.db.Begin()

	err := r.db.Create(newMerchant).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
