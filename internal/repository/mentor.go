package repository

import (
	"intern-bcc/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMentorRepository interface{
	GetMentor(mentor *domain.Mentors, mentorParam domain.Mentors) error 
	CreateMentor(newMentor *domain.Mentors) error
	UpdateMentor(mentor *domain.MentorUpdate, mentorId uuid.UUID) error
}

type MentorRepository struct {
	db *gorm.DB
}

func NewMentorRepository(db *gorm.DB) IMentorRepository {
	return &MentorRepository{db}
}

func(r *MentorRepository) GetMentor(mentor *domain.Mentors, mentorParam domain.Mentors) error {
	err := r.db.First(mentor, mentorParam).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *MentorRepository) CreateMentor(newMentor *domain.Mentors) error {
	err := r.db.Create(newMentor).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *MentorRepository) UpdateMentor(mentor *domain.MentorUpdate, mentorId uuid.UUID) error {
	err := r.db.Where("id = ?", mentorId).Updates(mentor).Error
	if err != nil{
		return err
	}

	return nil
}