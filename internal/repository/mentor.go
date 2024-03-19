package repository

import (
	"intern-bcc/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMentorRepository interface {
	GetMentor(mentor *domain.Mentors, mentorParam domain.MentorParam) error
	GetMentors(mentors *[]domain.Mentors) error
	CreateMentor(newMentor *domain.Mentors) error
	UpdateMentor(mentor *domain.MentorUpdate, mentorId uuid.UUID) error
}

type MentorRepository struct {
	db *gorm.DB
}

func NewMentorRepository(db *gorm.DB) IMentorRepository {
	return &MentorRepository{db}
}

func (r *MentorRepository) GetMentor(mentor *domain.Mentors, mentorParam domain.MentorParam) error {
	err := r.db.Preload("Experiences", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}, func(db *gorm.DB) *gorm.DB {
		return db.Limit(3)
	}).First(mentor, mentorParam).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *MentorRepository) GetMentors(mentors *[]domain.Mentors) error {
	err := r.db.Limit(15).Order("created_at desc").Find(mentors).Error
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
	err := r.db.Model(domain.Mentors{}).Where("id = ?", mentorId).Updates(mentor).Error
	if err != nil {
		return err
	}

	return nil
}
