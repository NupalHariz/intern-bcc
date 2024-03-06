package usecase

import (
	"intern-bcc/internal/repository"
)

type Usecase struct {
	UserUsecase IUserUsecase
}

type InitParam struct {
	Repository *repository.Repository
}

func NewUsecase(param InitParam) *Usecase {
	userUsecase := NewUserUsecase(param.Repository)

	return &Usecase{
		UserUsecase: userUsecase,
	}
}