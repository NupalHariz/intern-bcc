package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	UserRepository IUserRepository
	MerchantRepository IMerchantRepository
}

func NewRepository(db *gorm.DB) *Repository {
	userRepository := NewUserRepository(db)
	merchantRepository := NewMerchantRepository(db)

	return &Repository{
		UserRepository: userRepository,
		MerchantRepository: merchantRepository,
	}
}
