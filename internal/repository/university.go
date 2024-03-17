package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IUniversityRepository interface{
	GetUniversity(university *domain.Universities, universityParam domain.Universities) error
	CreateUniversity(university *domain.Universities) error
}

type UniversityRepository struct {
	db *gorm.DB
}

func NewUniversityRepository(db *gorm.DB) IUniversityRepository {
	return &UniversityRepository{db}
}

func (r *UniversityRepository) GetUniversity(university *domain.Universities, universityParam domain.Universities) error {
	err := r.db.First(university, universityParam).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UniversityRepository) CreateUniversity(university *domain.Universities) error{
	err := r.db.Create(university).Error
	if err != nil {
		return err
	}

	return nil
}