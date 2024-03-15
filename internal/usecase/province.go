package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type IProvinceUsecase interface{
	CreateProvince(provinceRequest domain.Province) any
}

type ProvinceUsecase struct {
	provinceRepository repository.IProvinceRepository
}

func NewProvinceUsecase(provinceRepository repository.IProvinceRepository) IProvinceUsecase {
	return &ProvinceUsecase{provinceRepository}
}

func (u *ProvinceUsecase) CreateProvince(provinceRequest domain.Province) any {
	var province domain.Province
	err := u.provinceRepository.GetProvince(&province, domain.Province{Province: provinceRequest.Province})
	if err == nil {
		return response.ErrorObject{
			Code: http.StatusBadRequest,
			Message: "province already exist",
			Err: err,
		}
	}

	province.Province = provinceRequest.Province
	err = u.provinceRepository.CreateProvince(&province)
	if err != nil {
		return response.ErrorObject{
			Code: http.StatusInternalServerError,
			Message: "an error occured when create province",
			Err: err,
		}
	}

	return nil
}