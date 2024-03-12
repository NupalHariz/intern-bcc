package repository

import (
	"fmt"
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IMerchantRepository interface {
	GetMerchant(merchant *domain.Merchants, param domain.MerchantParam) error
	CreateMerchant(newMerchant *domain.Merchants) error
	UpdateMerchant(updateMerchant *domain.Merchants) error
}

type MerchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) IMerchantRepository {
	return &MerchantRepository{db}
}

func (r *MerchantRepository) GetMerchant(merchant *domain.Merchants, param domain.MerchantParam) error {
	err := r.db.First(merchant, param).Error
	if err != nil {
		return err
	}

	return nil
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

func (r *MerchantRepository) UpdateMerchant(updateMerchant *domain.Merchants) error {
	tx := r.db.Begin()

	err := r.db.Where("id = ?", updateMerchant.Id).Updates(updateMerchant).Error
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
