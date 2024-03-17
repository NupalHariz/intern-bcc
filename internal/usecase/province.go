package usecase

import (
	"errors"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type IProvinceUsecase interface {
	CreateProvince(provinceRequest domain.Province) error
}

type ProvinceUsecase struct {
	provinceRepository repository.IProvinceRepository
}

func NewProvinceUsecase(provinceRepository repository.IProvinceRepository) IProvinceUsecase {
	return &ProvinceUsecase{provinceRepository}
}

func (u *ProvinceUsecase) CreateProvince(provinceRequest domain.Province) error {
	var province domain.Province
	err := u.provinceRepository.GetProvince(&province, domain.Province{Province: provinceRequest.Province})
	if err == nil {
		return response.NewError(http.StatusBadRequest, "province already exist", errors.New("can not make the same province"))
	}

	province.Province = provinceRequest.Province
	err = u.provinceRepository.CreateProvince(&province)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when create province", err)
	}

	return nil
}
