package usecase

import (
	"context"
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
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IMerchantUsecase interface {
	CreateMerchant(c *gin.Context, merchantRequest domain.MerchantRequest) any
	SendOtp(c *gin.Context, ctx context.Context) any
	VerifyOtp(c *gin.Context, ctx context.Context, verifyOtp domain.MerchantVerify) any
	UpdateMerchant(c *gin.Context, merchantId int, updateMerchant domain.UpdateMerchant) (domain.Merchants, any)
	UploadMerchantPhoto(c *gin.Context, merchantId int, merchantPhoto *multipart.FileHeader) (domain.Merchants, any)
}

type MerchantUsecase struct {
	redis                repository.IRedis
	merchantRepository   repository.IMerchantRepository
	provinceRepository   repository.IProvinceRepository
	universityRepository repository.IUniversityRepository
	jwt                  jwt.IJwt
	goMail               gomail.IGoMail
	supabase             supabase.ISupabase
}

func NewMerchantUsecase(merchantRepository repository.IMerchantRepository, redis repository.IRedis,
	jwt jwt.IJwt, goMail gomail.IGoMail, supabase supabase.ISupabase,
	universityRepository repository.IUniversityRepository, provinceRepository repository.IProvinceRepository) IMerchantUsecase {
	return &MerchantUsecase{
		redis:                redis,
		merchantRepository:   merchantRepository,
		provinceRepository:   provinceRepository,
		universityRepository: universityRepository,
		jwt:                  jwt,
		goMail:               goMail,
		supabase:             supabase,
	}
}

func (u *MerchantUsecase) CreateMerchant(c *gin.Context, merchantRequest domain.MerchantRequest) any {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}

	if merchantRequest.MerchantName == "" {
		merchantRequest.MerchantName = strings.Split(user.Email, "@")[0] + " Store's"
	}

	var province domain.Province
	err = u.provinceRepository.GetProvince(&province, domain.Province{Province: merchantRequest.Province})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "province does not exist",
			Err:     err,
		}
	}

	var university domain.Universities
	err = u.universityRepository.GetUniversity(&university, domain.Universities{University: merchantRequest.University})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "university does not exist",
			Err:     err,
		}
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err == nil {
		merchant := domain.Merchants{
			Id:           merchant.Id,
			MerchantName: merchantRequest.MerchantName,
			UniversityId: university.Id,
			Faculty:      merchantRequest.Faculty,
			ProvinceId:   province.Id,
			City:         merchantRequest.City,
			PhoneNumber:  merchantRequest.PhoneNumber,
			Instagram:    merchantRequest.Instagram,
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
		Id:           uuid.New(),
		UserId:       user.Id,
		MerchantName: merchantRequest.MerchantName,
		UniversityId: university.Id,
		Faculty:      merchantRequest.Faculty,
		ProvinceId:   province.Id,
		City:         merchantRequest.City,
		PhoneNumber:  merchantRequest.PhoneNumber,
		Instagram:    merchantRequest.Instagram,
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
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user",
			Err:     err,
		}
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "please create your merchant before verify",
			Err:     err,
		}
	}

	otp := rand.Intn(999999-100000) + 100000
	otpString := strconv.Itoa(otp)

	err = u.redis.SetOTP(ctx, user.Id, otpString)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when save otp to database",
			Err:     err,
		}
	}

	subject := "Verify Merchant Code"
	htmlBody := `<html>
	<p>Berikut adalah kode otp mu <strong>` + otpString + `</strong></p>
	</html>`

	err = u.goMail.SendGoMail(subject, htmlBody, user.Email)
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
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}

	stringOtp, err := u.redis.GetOTP(ctx, user.Id)
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

func (u *MerchantUsecase) UpdateMerchant(c *gin.Context, merchantId int, updateMerchant domain.UpdateMerchant) (domain.Merchants, any) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Merchants{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user",
			Err:     err,
		}
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{Id: merchantId})
	if err != nil {
		return domain.Merchants{}, response.ErrorObject{
			Code:    http.StatusFound,
			Message: "failed to get merchant",
			Err:     err,
		}
	}

	if user.Id != merchant.UserId {
		return domain.Merchants{}, response.ErrorObject{
			Code:    http.StatusUnauthorized,
			Message: "access denied",
			Err:     errors.New("can not edit other people merchant"),
		}
	}

	if updateMerchant.MerchantName != "" {
		merchant.MerchantName = updateMerchant.MerchantName
	}
	if updateMerchant.Instagram != "" {
		merchant.Instagram = updateMerchant.Instagram
	}
	if updateMerchant.PhoneNumber != "" {
		merchant.PhoneNumber = updateMerchant.PhoneNumber
	}

	err = u.merchantRepository.UpdateMerchant(&merchant)
	if err != nil {
		return domain.Merchants{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when update merchant",
			Err:     err,
		}
	}

	return merchant, nil
}

func (u *MerchantUsecase) UploadMerchantPhoto(c *gin.Context, merchantId int, merchantPhoto *multipart.FileHeader) (domain.Merchants, any) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Merchants{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user",
			Err:     err,
		}
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{Id: merchantId})
	if err != nil {
		return domain.Merchants{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get merchant",
			Err:     err,
		}
	}

	if user.Id != merchant.UserId {
		return domain.Merchants{}, response.ErrorObject{
			Code:    http.StatusUnauthorized,
			Message: "access denied",
			Err:     errors.New("can not edit other people merchant"),
		}
	}

	if merchant.MerchantPhoto != "" {
		err = u.supabase.Delete(merchant.MerchantPhoto)
		if err != nil {
			return domain.Merchants{}, response.ErrorObject{
				Code:    http.StatusInternalServerError,
				Message: "error occured when deleting old merchant photo",
				Err:     err,
			}
		}
	}

	merchantPhoto.Filename = fmt.Sprintf("%v-%v", time.Now().String(), merchantPhoto.Filename)
	if strings.Contains(merchantPhoto.Filename, " ") {
		merchantPhoto.Filename = strings.Replace(merchantPhoto.Filename, " ", "-", -1)
	}

	newMerchantPhoto, err := u.supabase.Upload(merchantPhoto)
	if err != nil {
		return domain.Merchants{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to upload photo",
			Err:     err,
		}
	}

	merchant.MerchantPhoto = newMerchantPhoto
	err = u.merchantRepository.UpdateMerchant(&merchant)
	if err != nil {
		return domain.Merchants{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when update merchant photo",
			Err:     err,
		}
	}

	return merchant, nil
}
