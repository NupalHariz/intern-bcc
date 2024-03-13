package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IMentorRepository interface{
	GetMentor(mentor *domain.Mentors, mentorParam domain.Mentors) error 
	CreateMentor(newMentor *domain.Mentors) error
	UpdateMentor (mentor *domain.Mentors) error
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
	tx := r.db.Begin()

	err := r.db.Create(newMentor).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *MentorRepository) UpdateMentor (mentor *domain.Mentors) error {
	tx := r.db.Begin()

	err := r.db.Where("id = ?", mentor.Id).Updates(mentor).Error
	if err != nil{
		tx.Rollback()
		return err
	}

	return nil
}