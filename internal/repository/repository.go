package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepository        IUserRepository
	ProductRepository     IProductRepository
	TransactionRepository ITransactionRepository
	MerchantSQLRepository IMerchantRepository
	RedisRepository       IRedis
	MentorRepository      IMentorRepository
	ExperienceRepository  IExperienceRepository
	CategoryRepository    ICategoryRepository
	InformationRepository IInformationRepository
}

func NewRepository(db *gorm.DB, r *redis.Client) *Repository {
	userRepository := NewUserRepository(db)
	productRepository := NewProductRepository(db)
	transactionRepository := NewTransactionRepository(db)
	merchantSQLRepository := NewMerchantRepository(db)
	redisRepository := NewRedis(r)
	mentorRepository := NewMentorRepository(db)
	experienceRepository := NewExperienceRepository(db)
	categoryRepository := NewCategoryRepository(db)
	informationRepository := NewInformationRepository(db)

	return &Repository{
		UserRepository:        userRepository,
		ProductRepository:     productRepository,
		TransactionRepository: transactionRepository,
		MerchantSQLRepository: merchantSQLRepository,
		RedisRepository:       redisRepository,
		MentorRepository:      mentorRepository,
		ExperienceRepository:  experienceRepository,
		CategoryRepository:    categoryRepository,
		InformationRepository: informationRepository,
	}
}
