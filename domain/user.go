package domain

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id        uuid.UUID `json:"-" gorm:"primary key"`
	Username  string    `json:"username" binding:"required"`
	Email     string    `json:"email" binding:"required,email" gorm:"unique"`
	Password  string    `json:"password" binding:"required"`
	IsAdmin   bool      `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type UserRequest struct {
	Username string `json:"Username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserParam struct {
	Id       uuid.UUID `json:"-"`
	Email    string    `json:"-"`
	Password string    `json:"-"`
}

type LoginResponse struct {
	JWT string `json:"jwt"`
}

//	ProfilePicture
