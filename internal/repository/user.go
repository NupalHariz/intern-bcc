package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUser(user *domain.Users, param domain.UserParam) error
	Register(newUser *domain.Users) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUser(user *domain.Users, param domain.UserParam) error {
	err := r.db.First(user, &param).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Register(newUser *domain.Users) error {
	tx := r.db.Begin()

	err := r.db.Create(newUser).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
