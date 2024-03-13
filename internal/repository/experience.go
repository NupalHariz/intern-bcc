package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IExperienceRepository interface {
	AddExperience(experience *domain.Experiences) error
}

type ExperienceRepository struct {
	db *gorm.DB
}

func NewExperienceRepository(db *gorm.DB) IExperienceRepository {
	return &ExperienceRepository{db}
}

func (r *ExperienceRepository) AddExperience(experience *domain.Experiences) error {
	tx := r.db.Begin()

	err := r.db.Create(experience).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
