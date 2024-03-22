package usecase

import (
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
	province := domain.Province{
		Province: provinceRequest.Province,
	}

	err := u.provinceRepository.CreateProvince(&province)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when create province", err)
	}

	return nil
}
