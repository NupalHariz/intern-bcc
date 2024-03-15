package usecase

import (
	"errors"
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/gomail"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"intern-bcc/pkg/supabase"
	"math/rand"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type IUserUsecase interface {
	GetUser(param domain.UserParam) (domain.Users, any)
	Register(userRequest domain.UserRequest) any
	Login(userLogin domain.UserLogin) (domain.LoginResponse, any)
	UpdateUser(c *gin.Context, userId uuid.UUID, userUpdate domain.UserUpdate) (domain.Users, any)
	UploadUserPhoto(c *gin.Context, userId uuid.UUID, userPhoto *multipart.FileHeader) any
	PasswordRecovery(userParam domain.UserParam, ctx context.Context) any
	ChangePassword(ctx context.Context, email string, verPass string, passwordRequest domain.PasswordUpdate) any
	LikeProduct(c *gin.Context, productId int) any
	DeleteLikeProduct(c *gin.Context, productId int) any
}

type UserUsecase struct {
	userRepository    repository.IUserRepository
	productRepository repository.IProductRepository
	jwt               jwt.IJwt
	supabase          supabase.ISupabase
	redis             repository.IRedis
	goMail            gomail.IGoMail
}

func NewUserUsecase(userRepository repository.IUserRepository, productRepository repository.IProductRepository, jwt jwt.IJwt,
	supabase supabase.ISupabase, redis repository.IRedis, goMail gomail.IGoMail) IUserUsecase {
	return &UserUsecase{
		userRepository:    userRepository,
		productRepository: productRepository,
		jwt:               jwt,
		supabase:          supabase,
		redis:             redis,
		goMail:            goMail,
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

func (u *UserUsecase) UpdateUser(c *gin.Context, userId uuid.UUID, userUpdate domain.UserUpdate) (domain.Users, any) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Users{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}

	if userId != user.Id {
		return domain.Users{}, response.ErrorObject{
			Code:    http.StatusUnauthorized,
			Message: "access denied",
			Err:     errors.New("can not change other people profile"),
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

func (u *UserUsecase) UploadUserPhoto(c *gin.Context, userId uuid.UUID, userPhoto *multipart.FileHeader) any {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user",
			Err:     err,
		}
	}

	if user.Id != userId {
		return response.ErrorObject{
			Code:    http.StatusUnauthorized,
			Message: "access denied",
			Err:     errors.New("can not change other people profile"),
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
			Message: "failed to upload photo",
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

func (u *UserUsecase) PasswordRecovery(userParam domain.UserParam, ctx context.Context) any {
	fmt.Println(userParam.Email)
	var user domain.Users
	err := u.userRepository.GetUser(&user, domain.UserParam{Email: userParam.Email})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "account not found",
			Err:     err,
		}
	}

	fmt.Println(user.Email)

	var emailVerPassword string
	emailVerPassword, emailVerHash, err := newEmailVerification()
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "unable to send recovery email",
			Err:     err,
		}
	}

	fmt.Println(emailVerHash)
	err = u.redis.SetEmailVerHash(ctx, user.Email, emailVerHash)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "unable to send recovery email",
			Err:     err,
		}
	}

	link := fmt.Sprintf("http://localhost:8080/accountrecovery/%v/%v", user.Email, emailVerPassword)
	err = u.goMail.SendGoMail(link, user.Email)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "unable to send recovery email",
			Err:     err,
		}
	}

	return nil
}

func (u *UserUsecase) ChangePassword(ctx context.Context, email string, verPass string, passwordRequest domain.PasswordUpdate) any {
	var user domain.Users
	err := u.userRepository.GetUser(&user, domain.UserParam{Email: email})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "user not found",
			Err:     err,
		}
	}

	verPassHash, err := u.redis.GetEmailVerHash(ctx, email)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to get ver password from database",
			Err:     err,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(verPassHash), []byte(verPass))
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "an error occured",
			Err:     err,
		}
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(passwordRequest.Password), 10)
	if err != nil {
		return response.ErrorObject{
			Code: http.StatusBadRequest,
			Message: "an error occured when update password",
			Err: err,
		}
	}

	user.Password = string(hashedPass)
	err = u.userRepository.UpdateUser(&user)
	if err != nil {
		return response.ErrorObject{
			Code: http.StatusInternalServerError,
			Message: "an error occured when update password",
			Err: err,
		}
	}

	return nil
}

func newEmailVerification() (string, string, error) {
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	emailVerPassword := string(emailVerRandRune)
	var emailVerPWhash []byte

	emailVerPWhash, err := bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	emailVerHash := string(emailVerPWhash)

	return emailVerPassword, emailVerHash, nil
}

func (u *UserUsecase) LikeProduct(c *gin.Context, productId int) any {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user",
			Err:     err,
		}
	}

	var product domain.Products
	err = u.productRepository.GetProduct(&product, &domain.ProductParam{Id: productId})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to find product",
			Err:     err,
		}
	}

	likeProduct := domain.LikeProduct{
		UserId:    user.Id,
		ProductId: product.Id,
	}

	err = u.userRepository.LikeProduct(&likeProduct)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when like product",
			Err:     err,
		}
	}

	return nil
}

func (u *UserUsecase) DeleteLikeProduct(c *gin.Context, productId int) any {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user",
			Err:     err,
		}
	}

	var likedProduct domain.LikeProduct
	err = u.userRepository.GetLikeProduct(&likedProduct, domain.LikeProduct{UserId: user.Id, ProductId: productId})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get liked product",
			Err:     err,
		}
	}

	err = u.userRepository.DeleteLikeProduct(&likedProduct)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to delete like product",
			Err:     err,
		}
	}

	return nil
}
