package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IProvinceRepository interface {
	GetProvince(province *domain.Province, provinceParam domain.Province) error
	CreateProvince(province *domain.Province) error
}

type ProvinceRepository struct {
	db *gorm.DB
}

func NewProvinceRepository(db *gorm.DB) IProvinceRepository {
	return &ProvinceRepository{db}
}

func (r *ProvinceRepository) GetProvince(province *domain.Province, provinceParam domain.Province) error {
	err := r.db.First(province, provinceParam).Error
	// if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return response.NewError(http.StatusNotFound, "province not found", err)
	// } else {
	// 	return response.NewError(http.StatusInternalServerError, "failed get province", err)
	// }
	if err != nil {
		return err
	}

	return nil
}

func (r *ProvinceRepository) CreateProvince(province *domain.Province) error {
	err := r.db.Create(province).Error
	if err != nil {
		return err
	}

	return nil
}
