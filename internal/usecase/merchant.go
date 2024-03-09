package usecase

import (
	"context"
	"errors"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/gomail"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type IMerchantUsecase interface {
	CreateMerchant(c *gin.Context, merchantRequest domain.MerchantRequest) any
	SendOtp(c *gin.Context, ctx context.Context) any
	VerifyOtp(c *gin.Context, ctx context.Context, verifyOtp domain.MerchantVerify) any
}

type MerchantUsecase struct {
	merchantRedis      repository.IMerchantRedis
	merchantRepository repository.IMerchantRepository
	userRepository repository.IUserRepository
}

func NewMerchantUsecase(merchantRepository repository.IMerchantRepository, merchantRedis repository.IMerchantRedis, userRepository repository.IUserRepository) IMerchantUsecase {
	return &MerchantUsecase{
		merchantRedis:      merchantRedis,
		merchantRepository: merchantRepository,
		userRepository: userRepository,
	}
}

func (u *MerchantUsecase) CreateMerchant(c *gin.Context, merchantRequest domain.MerchantRequest) any {
	userId, err := jwt.GetLoginUserId(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}
	var user domain.Users
	err = u.userRepository.GetUser(&user, domain.UserParam{Id: userId})
	if err != nil {
		return response.ErrorObject {
			Code: http.StatusNotFound,
			Message: "failed to get user",
			Err: err,
		}
	}

	if merchantRequest.StoreName == "" {
		merchantRequest.StoreName = strings.Split(user.Email, "@")[0] + " Store's"
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err == nil {
		merchant := domain.Merchants{
			Id:          merchant.Id,
			StoreName:   merchantRequest.StoreName,
			University:  merchantRequest.University,
			Faculty:     merchantRequest.Faculty,
			Province:    merchantRequest.Province,
			City:        merchantRequest.City,
			PhoneNumber: merchantRequest.PhoneNumber,
			Instagram:   merchantRequest.Instagram,
		}

		err = u.merchantRepository.UpdateMerchant(&merchant)
		if err != nil {
			return response.ErrorObject{
				Code:    http.StatusInternalServerError,
				Message: "error occured when update account",
				Err:     err,
			}
		}

		return nil
	}

	newMerchant := domain.Merchants{
		UserId:      user.Id,
		StoreName:   merchantRequest.StoreName,
		University:  merchantRequest.University,
		Faculty:     merchantRequest.Faculty,
		Province:    merchantRequest.Province,
		City:        merchantRequest.City,
		PhoneNumber: merchantRequest.PhoneNumber,
		Instagram:   merchantRequest.Instagram,
	}

	err = u.merchantRepository.CreateMerchant(&newMerchant)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when create account",
			Err:     err,
		}
	}

	return nil
}

func (u *MerchantUsecase) SendOtp(c *gin.Context, ctx context.Context) any {
	userId, err := jwt.GetLoginUserId(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}
	var user domain.Users
	err = u.userRepository.GetUser(&user, domain.UserParam{Id: userId})
	if err != nil {
		return response.ErrorObject {
			Code: http.StatusNotFound,
			Message: "failed to get user",
			Err: err,
		}
	}

	otp := rand.Intn(999999-100000) + 100000
	otpString := strconv.Itoa(otp)

	err = u.merchantRedis.SetOTP(ctx, user.Id, otpString)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when save otp to database",
			Err:     err,
		}
	}

	err = gomail.SendGoMail(otpString, user.Email)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to send otp email",
			Err:     err,
		}
	}

	return nil
}

func (u *MerchantUsecase) VerifyOtp(c *gin.Context, ctx context.Context, verifyOtp domain.MerchantVerify) any {
	userId, err := jwt.GetLoginUserId(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}
	var user domain.Users
	err = u.userRepository.GetUser(&user, domain.UserParam{Id: userId})
	if err != nil {
		return response.ErrorObject {
			Code: http.StatusNotFound,
			Message: "failed to get user",
			Err: err,
		}
	}

	stringOtp, err := u.merchantRedis.GetOTP(ctx, user.Id)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when getting otp from redis",
			Err:     err,
		}
	}

	if verifyOtp.VerifyOtp != stringOtp {
		return response.ErrorObject{
			Code:    http.StatusUnauthorized,
			Message: "invalid token",
			Err:     errors.New("wrong token"),
		}
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "an error occured when getting merchant",
			Err:     err,
		}
	}

	merchant.IsActive = true

	err = u.merchantRepository.UpdateMerchant(&merchant)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "error occured when update account",
			Err:     err,
		}
	}

	return nil
}
