package usecase

import (
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"intern-bcc/pkg/supabase"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	GetUser(param domain.UserParam) (domain.Users, any)
	Register(userRequest domain.UserRequest) any
	Login(userLogin domain.UserLogin) (domain.LoginResponse, any)
	UpdateUser(c *gin.Context, userUpdate domain.UserUpdate) (domain.Users, any)
	UploadPhoto(c *gin.Context, userPhoto *multipart.FileHeader) any
}

type UserUsecase struct {
	userRepository repository.IUserRepository
	jwt            jwt.IJwt
	supabase       supabase.ISupabase
}

func NewUserUsecase(userRepository repository.IUserRepository, jwt jwt.IJwt, supabase supabase.ISupabase) IUserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
		jwt:            jwt,
		supabase:       supabase,
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
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		IsAdmin:  userRequest.IsAdmin,
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

	tokenString, err := u.jwt.GenerateToken(user.Id)
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

func (u *UserUsecase) UpdateUser(c *gin.Context, userUpdate domain.UserUpdate) (domain.Users, any) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Users{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}

	user = checkNullUpdateUser(user, userUpdate)

	err = u.userRepository.UpdateUser(&user)
	if err != nil {
		return domain.Users{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "error occured when update user",
			Err:     err,
		}
	}

	return user, nil
}

func (u *UserUsecase) UploadPhoto(c *gin.Context, userPhoto *multipart.FileHeader) any {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		fmt.Println("148", err)

		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}

	if user.ProfilePicture != "" {
		err = u.supabase.Delete(user.ProfilePicture)
		if err != nil {

			return response.ErrorObject{
				Code:    http.StatusInternalServerError,
				Message: "error occured when deleting old profile picture",
				Err:     err,
			}
		}
	}

	newProfilePicture, err := u.supabase.Upload(userPhoto)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to uplaod photo",
			Err:     err,
		}
	}

	user.ProfilePicture = newProfilePicture
	err = u.userRepository.UpdateUser(&user)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "error occured when update user",
			Err:     err,
		}
	}

	return nil
}

func checkNullUpdateUser(user domain.Users, userUpdate domain.UserUpdate) domain.Users {
	if userUpdate.Gender != "" {
		user.Gender = userUpdate.Gender
	}

	if userUpdate.Name != "" {
		user.Name = userUpdate.Name
	}

	if userUpdate.PlaceBirth != "" {
		user.PlaceBirth = userUpdate.PlaceBirth
	}

	if userUpdate.DateBirth != "" {
		user.DateBirth = userUpdate.DateBirth
	}

	return user
}
