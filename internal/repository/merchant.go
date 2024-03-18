package repository

import (
	"intern-bcc/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMerchantRepository interface {
	GetMerchant(merchant *domain.Merchants, param domain.MerchantParam) error
	CreateMerchant(newMerchant *domain.Merchants) error
	UpdateMerchant(updateMerchant *domain.UpdateMerchant, merchantId uuid.UUID) error
}

type MerchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) IMerchantRepository {
	return &MerchantRepository{db}
}

func (r *MerchantRepository) GetMerchant(merchant *domain.Merchants, param domain.MerchantParam) error {
	err := r.db.Preload("University").Preload("Province").First(merchant, param).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *MerchantRepository) CreateMerchant(newMerchant *domain.Merchants) error {
	err := r.db.Create(newMerchant).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *MerchantRepository) UpdateMerchant(updateMerchant *domain.UpdateMerchant, merchantId uuid.UUID) error {
	err := r.db.Debug().Where("id = ?", merchantId).Updates(updateMerchant).Error
	if err != nil {
		return err
	}

	return nil
}
