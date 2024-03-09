package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	GetUser(param domain.UserParam) (domain.Users, any)
	Register(userRequest domain.UserRequest) any
	Login(userLogin domain.UserLogin) (domain.LoginResponse, any)
}

type UserUsecase struct {
	userRepository repository.IUserRepository
}

func NewUserUsecase(userRepository repository.IUserRepository) IUserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (u *UserUsecase) GetUser(param domain.UserParam) (domain.Users, any) {
	var user domain.Users
	err := u.userRepository.GetUser(&user, param)
	if err != nil {
		return user, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user",
			Err:     err,
		}
	}

	return user, nil
}

func (u *UserUsecase) Register(userRequest domain.UserRequest) any {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 10)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Err:     err,
			Message: "error when hashing password",
		}
	}

	NewUser := domain.Users{
		Id:       uuid.New(),
		Username: userRequest.Username,
		Email:    userRequest.Email,
		Password: string(hashPassword),
	}

	err = u.userRepository.Register(&NewUser)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Err:     err,
			Message: "error occured when creating user",
		}
	}

	return nil
}

func (u *UserUsecase) Login(userLogin domain.UserLogin) (domain.LoginResponse, any) {
	var user domain.Users
	err := u.userRepository.GetUser(&user, domain.UserParam{
		Email: userLogin.Email,
	})
	if err != nil {
		return domain.LoginResponse{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "email or invalid",
			Err:     err,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		return domain.LoginResponse{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: " or password invalid",
			Err:     err,
		}
	}

	tokenString, err := jwt.GenerateToken(user.Id)
	if err != nil {
		return domain.LoginResponse{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to generate jwt token",
			Err:     err,
		}
	}

	loginUser := domain.LoginResponse{
		JWT: tokenString,
	}

	return loginUser, nil
}
