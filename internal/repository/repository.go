package repository

import (
	"intern-bcc/pkg/redis"

	"gorm.io/gorm"
)

const (
	KeySetOtp                = "otp:set:id:%v"
	KeySetPasswordRecovery   = "recovery:set:name:%v"
	KeySetInformationNmentor = "get:all:%v"
	KeySetProducts           = "get:all:product:%v"
	Limit                    = 6
)

type Repository struct {
	UserRepository        IUserRepository
	ProductRepository     IProductRepository
	TransactionRepository ITransactionRepository
	MerchantSQLRepository IMerchantRepository
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
	userRepository := NewUserRepository(db, repositoryParam.Redis)
	productRepository := NewProductRepository(db, repositoryParam.Redis)
	transactionRepository := NewTransactionRepository(db)
	merchantSQLRepository := NewMerchantRepository(db, repositoryParam.Redis)
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
		MentorRepository:      mentorRepository,
		ExperienceRepository:  experienceRepository,
		CategoryRepository:    categoryRepository,
		InformationRepository: informationRepository,
		UniversityRepository:  universityRepository,
		ProvinceRepository:    provinceRepository,
	}
}
