package repository

import (
	"gorm.io/gorm"
	"intern-bcc/pkg/redis"
)

type Repository struct {
	UserRepository        IUserRepository
	ProductRepository     IProductRepository
	TransactionRepository ITransactionRepository
	MerchantSQLRepository IMerchantRepository
	// RedisRepository       IRedis
	MentorRepository      IMentorRepository
	ExperienceRepository  IExperienceRepository
	CategoryRepository    ICategoryRepository
	InformationRepository IInformationRepository
	UniversityRepository  IUniversityRepository
	ProvinceRepository    IProvinceRepository
}

type RepositoryParam struct {
	Redis redis.IRedis
}

func NewRepository(db *gorm.DB, repositoryParam RepositoryParam) *Repository {
	userRepository := NewUserRepository(db)
	productRepository := NewProductRepository(db, repositoryParam.Redis)
	transactionRepository := NewTransactionRepository(db)
	merchantSQLRepository := NewMerchantRepository(db)
	// redisRepository := NewRedis(r)
	mentorRepository := NewMentorRepository(db, repositoryParam.Redis)
	experienceRepository := NewExperienceRepository(db)
	categoryRepository := NewCategoryRepository(db)
	informationRepository := NewInformationRepository(db, repositoryParam.Redis)
	universityRepository := NewUniversityRepository(db)
	provinceRepository := NewProvinceRepository(db)

	return &Repository{
		UserRepository:        userRepository,
		ProductRepository:     productRepository,
		TransactionRepository: transactionRepository,
		MerchantSQLRepository: merchantSQLRepository,
		// RedisRepository:       redisRepository,
		MentorRepository:      mentorRepository,
		ExperienceRepository:  experienceRepository,
		CategoryRepository:    categoryRepository,
		InformationRepository: informationRepository,
		UniversityRepository:  universityRepository,
		ProvinceRepository:    provinceRepository,
	}
}
