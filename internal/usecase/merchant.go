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
	GetMerchant(c *gin.Context) (domain.MerchantProfileResponse, error)
	CreateMerchant(c *gin.Context, merchantRequest domain.MerchantRequest) error
	SendOtp(c *gin.Context, ctx context.Context) error
	VerifyOtp(c *gin.Context, ctx context.Context, verifyOtp domain.MerchantVerify) error
	UpdateMerchant(c *gin.Context, merchantId uuid.UUID, updateMerchant domain.UpdateMerchant) (domain.MerchantProfileResponse, error)
	UploadMerchantPhoto(c *gin.Context, merchantId uuid.UUID, merchantPhoto *multipart.FileHeader) (domain.MerchantProfileResponse, error)
}

type MerchantUsecase struct {
	merchantRepository   repository.IMerchantRepository
	provinceRepository   repository.IProvinceRepository
	universityRepository repository.IUniversityRepository
	jwt                  jwt.IJwt
	goMail               gomail.IGoMail
	supabase             supabase.ISupabase
}

func NewMerchantUsecase(merchantRepository repository.IMerchantRepository,
	jwt jwt.IJwt, goMail gomail.IGoMail, supabase supabase.ISupabase,
	universityRepository repository.IUniversityRepository, provinceRepository repository.IProvinceRepository) IMerchantUsecase {
	return &MerchantUsecase{
		merchantRepository:   merchantRepository,
		provinceRepository:   provinceRepository,
		universityRepository: universityRepository,
		jwt:                  jwt,
		goMail:               goMail,
		supabase:             supabase,
	}
}

func (u *MerchantUsecase) GetMerchant(c *gin.Context) (domain.MerchantProfileResponse, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusInternalServerError, "an error occured when get merchant", err)
	}

	merchantResponse := domain.MerchantProfileResponse{
		Id:            merchant.Id,
		MerchantName:  merchant.MerchantName,
		Province:      merchant.Province.Province,
		City:          merchant.City,
		University:    merchant.University.University,
		Faculty:       merchant.Faculty,
		PhoneNumber:   merchant.PhoneNumber,
		Instagram:     merchant.Instagram,
		MerchantPhoto: merchant.MerchantPhoto,
	}

	return merchantResponse, nil
}

func (u *MerchantUsecase) CreateMerchant(c *gin.Context, merchantRequest domain.MerchantRequest) error {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	if merchantRequest.MerchantName == "" {
		merchantRequest.MerchantName = strings.Split(user.Email, "@")[0] + " Store's"
	}

	var province domain.Province
	err = u.provinceRepository.GetProvince(&province, domain.Province{Province: merchantRequest.Province})
	if err != nil {
		return response.NewError(http.StatusBadRequest, "province does not exist", err)
	}

	var university domain.Universities
	err = u.universityRepository.GetUniversity(&university, domain.Universities{University: merchantRequest.University})
	if err != nil {
		return response.NewError(http.StatusBadRequest, "university does not exist", err)
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})

	if merchant.IsActive {
		return response.NewError(http.StatusBadRequest, "an error occured when create merchant", errors.New("you already have merchant"))
	}
	if err == nil {
		updateMerchant := domain.UpdateMerchant{
			MerchantName: merchantRequest.MerchantName,
			UniversityId: university.Id,
			Faculty:      merchantRequest.Faculty,
			ProvinceId:   province.Id,
			City:         merchantRequest.City,
			PhoneNumber:  merchantRequest.PhoneNumber,
			Instagram:    merchantRequest.Instagram,
		}

		err = u.merchantRepository.UpdateMerchant(&updateMerchant, merchant.Id)
		if err != nil {
			return response.NewError(http.StatusInternalServerError, "an error occured when update product", err)
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
		return response.NewError(http.StatusInternalServerError, "an errpr occured when create product", err)
	}

	return nil
}

func (u *MerchantUsecase) SendOtp(c *gin.Context, ctx context.Context) error {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return response.NewError(http.StatusNotFound, "please create your merchant before verify", err)
	}

	otp := rand.Intn(999999-100000) + 100000
	otpString := strconv.Itoa(otp)

	err = u.merchantRepository.CreateOTP(ctx, user.Id, otpString)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when make otp", err)
	}

	subject := "Verify Merchant Code"
	htmlBody := `<html>
	<p>Berikut adalah kode otp mu <strong>` + otpString + `</strong></p>
	</html>`

	err = u.goMail.SendGoMail(subject, htmlBody, user.Email)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "failed to send otp email", err)
	}

	return nil
}

func (u *MerchantUsecase) VerifyOtp(c *gin.Context, ctx context.Context, verifyOtp domain.MerchantVerify) error {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	stringOtp, err := u.merchantRepository.GetOTP(ctx, user.Id)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when get otp", err)
	}

	if verifyOtp.VerifyOtp != stringOtp {
		return response.NewError(http.StatusUnauthorized, "invalid token", errors.New("wrong token"))
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get merchant", err)
	}

	merchant.IsActive = true

	err = u.merchantRepository.UpdateMerchant(&domain.UpdateMerchant{IsActive: true}, merchant.Id)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when update account", err)
	}

	return nil
}

func (u *MerchantUsecase) UpdateMerchant(c *gin.Context, merchantId uuid.UUID, updateMerchant domain.UpdateMerchant) (domain.MerchantProfileResponse, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{Id: merchantId})
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusNotFound, "an error occured when get merchant", err)
	}

	if user.Id != merchant.UserId {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusUnauthorized, "access denied", errors.New("can not edit other people merchant"))
	}

	err = u.merchantRepository.UpdateMerchant(&updateMerchant, merchant.Id)
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusInternalServerError, "an error occured when update merchant", err)
	}

	var updatedMerchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&updatedMerchant, domain.MerchantParam{Id: merchant.Id})
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated merchant", err)
	}

	updatedMerchantResponse := domain.MerchantProfileResponse{
		Id:            updatedMerchant.Id,
		MerchantName:  updatedMerchant.MerchantName,
		Province:      updatedMerchant.Province.Province,
		City:          updatedMerchant.City,
		University:    updatedMerchant.University.University,
		Faculty:       updatedMerchant.Faculty,
		PhoneNumber:   updatedMerchant.PhoneNumber,
		Instagram:     updatedMerchant.Instagram,
		MerchantPhoto: updatedMerchant.MerchantPhoto,
	}

	return updatedMerchantResponse, nil
}

func (u *MerchantUsecase) UploadMerchantPhoto(c *gin.Context, merchantId uuid.UUID, merchantPhoto *multipart.FileHeader) (domain.MerchantProfileResponse, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{Id: merchantId})
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusNotFound, "an error occured when get merchant", err)
	}

	if user.Id != merchant.UserId {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusUnauthorized, "access denied", errors.New("can not edit other people merchant"))
	}

	if merchant.MerchantPhoto != "" {
		err = u.supabase.Delete(merchant.MerchantPhoto)
		if err != nil {
			return domain.MerchantProfileResponse{}, response.NewError(http.StatusInternalServerError, "an error occured when delete old merchant photo", err)
		}
	}

	merchantPhoto.Filename = fmt.Sprintf("%v-%v", time.Now().String(), merchantPhoto.Filename)
	merchantPhoto.Filename = strings.Replace(merchantPhoto.Filename, " ", "-", -1)

	newMerchantPhoto, err := u.supabase.Upload(merchantPhoto)
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusInternalServerError, "an error occured when upload photo", err)
	}

	err = u.merchantRepository.UpdateMerchant(&domain.UpdateMerchant{MerchantPhoto: newMerchantPhoto}, merchant.Id)
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusInternalServerError, "an error occured when update merchant photo", err)
	}

	var updatedMerchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&updatedMerchant, domain.MerchantParam{Id: merchant.Id})
	if err != nil {
		return domain.MerchantProfileResponse{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated merchant", err)
	}

	updatedMerchantResponse := domain.MerchantProfileResponse{
		Id:            updatedMerchant.Id,
		MerchantName:  updatedMerchant.MerchantName,
		Province:      updatedMerchant.Province.Province,
		City:          updatedMerchant.City,
		University:    updatedMerchant.University.University,
		Faculty:       updatedMerchant.Faculty,
		PhoneNumber:   updatedMerchant.PhoneNumber,
		Instagram:     updatedMerchant.Instagram,
		MerchantPhoto: updatedMerchant.MerchantPhoto,
	}

	return updatedMerchantResponse, nil
}
