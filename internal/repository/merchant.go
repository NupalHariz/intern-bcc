package repository

import (
	"context"
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/pkg/redis"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMerchantRepository interface {
	GetMerchant(merchant *domain.Merchants, param domain.MerchantParam) error
	CreateMerchant(newMerchant *domain.Merchants) error
	UpdateMerchant(updateMerchant *domain.UpdateMerchant, merchantId uuid.UUID) error
	CreateOTP(ctx context.Context, id uuid.UUID, otp string) error
	GetOTP(ctx context.Context, userId uuid.UUID) (string, error)
}

type MerchantRepository struct {
	db    *gorm.DB
	redis redis.IRedis
}

func NewMerchantRepository(db *gorm.DB, redis redis.IRedis) IMerchantRepository {
	return &MerchantRepository{db, redis}
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
	err := r.db.Model(domain.Merchants{}).Where("id = ?", merchantId).Updates(updateMerchant).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *MerchantRepository) CreateOTP(ctx context.Context, userId uuid.UUID, otp string) error{
	key := fmt.Sprintf(KeySetOtp, userId)
	err := r.redis.SetRedis(ctx, key, otp, 2*time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func (r *MerchantRepository) GetOTP(ctx context.Context, userId uuid.UUID) (string, error) {
	key := fmt.Sprintf(KeySetOtp, userId)
	otpString, err := r.redis.GetRedis(ctx, key)
	if err != nil {
		return "", err
	}

	return otpString, nil
}
