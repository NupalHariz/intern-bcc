package usecase

import (
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/gomail"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/midtrans"
	"intern-bcc/pkg/supabase"
)

type Usecase struct {
	UserUsecase        IUserUsecase
	TransactionUsecase ITransactionUsecase
	ProductUsecase     IProductUsecase
	MentorUsecase      IMentorUsecase
	MerchantUsecase    IMerchantUsecase
	ExperienceUsecase  IExperieceUsecase
	CategoryUsecase    ICategoryUsecase
	InformationUsecase IInformationUsecase
}

type UsecaseParam struct {
	Repository *repository.Repository
	Jwt        jwt.IJwt
	Supabase   supabase.ISupabase
	Midtrans   midtrans.IMidTrans
	GoMail     gomail.IGoMail
}

func NewUsecase(usecaseParam UsecaseParam) *Usecase {
	userUsecase := NewUserUsecase(usecaseParam.Repository.UserRepository, usecaseParam.Repository.ProductRepository, usecaseParam.Jwt, usecaseParam.Supabase)
	transactionUsecase := NewTransactionUsecase(usecaseParam.Repository.TransactionRepository, usecaseParam.Repository.UserRepository , usecaseParam.Jwt, usecaseParam.Midtrans)
	productUsecase := NewProductUsecase(usecaseParam.Repository.ProductRepository, usecaseParam.Jwt, usecaseParam.Repository.MerchantSQLRepository, usecaseParam.Repository.CategoryRepository,usecaseParam.Supabase)
	mentorUsecase := NewMentorUsecase(usecaseParam.Repository.MentorRepository, usecaseParam.Jwt, usecaseParam.Supabase)
	merchantUsecase := NewMerchantUsecase(usecaseParam.Repository.MerchantSQLRepository, usecaseParam.Repository.MerchantRedisRepository, usecaseParam.Jwt, usecaseParam.GoMail, usecaseParam.Supabase)
	experienceUsecase := NewExperienceRepository(usecaseParam.Repository.ExperienceRepository)
	categoryUsecase := NewCategoryUsecase(usecaseParam.Repository.CategoryRepository)
	informationUsecase := NewInformatinUsecase(usecaseParam.Repository.InformationRepository, usecaseParam.Repository.CategoryRepository, usecaseParam.Supabase)

	return &Usecase{
		UserUsecase:        userUsecase,
		TransactionUsecase: transactionUsecase,
		ProductUsecase:     productUsecase,
		MentorUsecase:      mentorUsecase,
		MerchantUsecase:    merchantUsecase,
		ExperienceUsecase:  experienceUsecase,
		CategoryUsecase: categoryUsecase,
		InformationUsecase: informationUsecase,
	}
}
