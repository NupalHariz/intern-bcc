package repository

import (
	"gorm.io/gorm"
)

type IUserRepository interface {
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}
