package usecase

import (
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
)

type Usecase struct {
	UserUsecase IUserUsecase
	MerchantUsecase IMerchantUsecase
}

type InitParam struct {
	Repository *repository.Repository
	JWT        jwt.IJwt
}

func NewUsecase(param InitParam) *Usecase {
	userUsecase := NewUserUsecase(param.Repository.UserRepository, param.JWT)
	merchantUsecase := NewMerchantUsecase(param.Repository.MerchantRepository, param.JWT)

	return &Usecase{
		UserUsecase: userUsecase,
		MerchantUsecase: merchantUsecase,
	}
}
