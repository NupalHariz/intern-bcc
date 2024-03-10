package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IMentorRepository interface{
	CreateMentor(newMentor *domain.Mentors) error
}

type MentorRepository struct {
	db *gorm.DB
}

func NewMentorRepository(db *gorm.DB) IMentorRepository {
	return &MentorRepository{db}
}

func (r *MentorRepository) CreateMentor(newMentor *domain.Mentors) error {
	tx := r.db.Begin()

	err := r.db.Create(newMentor).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}