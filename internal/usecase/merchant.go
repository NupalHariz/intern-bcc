package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type IMerchantUsecase interface {
	CreateMerchant(c *gin.Context, merchantRequest domain.MerchantRequest) any
}

type MerchantUsecase struct {
	merchantRepository repository.IMerchantRepository
	jwtAuth            jwt.IJwt
}

func NewMerchantUsecase(merchantRepository repository.IMerchantRepository, jwtAuth jwt.IJwt) IMerchantUsecase {
	return &MerchantUsecase{
		merchantRepository: merchantRepository,
		jwtAuth:            jwtAuth,
	}
}

func (u *MerchantUsecase) CreateMerchant(c *gin.Context, merchantRequest domain.MerchantRequest) any {
	var user domain.Users
	user, err := u.jwtAuth.GetLoginUser(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get account",
			Err:     err,
		}
	}

	if merchantRequest.StoreName == "" {
		merchantRequest.StoreName = strings.Split(user.Email, "@")[0] + " Store's"
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
