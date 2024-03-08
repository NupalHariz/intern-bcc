package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepository     IUserRepository
	MerchantRepository IMerchantRepository
	MerchantRedis      IMerchantRedis
}

func NewRepository(db *gorm.DB, r *redis.Client) *Repository {
	userRepository := NewUserRepository(db)
	merchantRepository := NewMerchantRepository(db)
	merchantRedisRepository := NewMerchantRedis(r)

	return &Repository{
		UserRepository:     userRepository,
		MerchantRepository: merchantRepository,
		MerchantRedis:      merchantRedisRepository,
	}
}
