package usecase

import "intern-bcc/internal/repository"

type IUserUsecase interface{}

type UserUsecase struct {
	userRepository repository.IUserRepository
}

func NewUserUsecase(userRepository repository.IUserRepository) *UserUsecase{
	return &UserUsecase{userRepository}
}
