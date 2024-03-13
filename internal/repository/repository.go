package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepository          IUserRepository
	ProductRepository       IProductRepository
	TransactionRepository   ITransactionRepository
	MerchantSQLRepository   IMerchantRepository
	MerchantRedisRepository IMerchantRedis
	MentorRepository        IMentorRepository
	ExperienceRepository    IExperienceRepository
	CategoryRepository      ICategoryRepository
	InformationRepository   IInformationRepository
}

func NewRepository(db *gorm.DB, r *redis.Client) *Repository {
	userRepository := NewUserRepository(db)
	productRepository := NewProductRepository(db)
	transactionRepository := NewTransactionRepository(db)
	merchantSQLRepository := NewMerchantRepository(db)
	merchantRedisRepository := NewMerchantRedis(r)
	mentorRepository := NewMentorRepository(db)
	experienceRepository := NewExperienceRepository(db)
	categoryRepository := NewCategoryRepository(db)
	informationRepository := NewInformationRepository(db)

	return &Repository{
		UserRepository:          userRepository,
		ProductRepository:       productRepository,
		TransactionRepository:   transactionRepository,
		MerchantSQLRepository:   merchantSQLRepository,
		MerchantRedisRepository: merchantRedisRepository,
		MentorRepository:        mentorRepository,
		ExperienceRepository:    experienceRepository,
		CategoryRepository:      categoryRepository,
		InformationRepository:   informationRepository,
	}
}
