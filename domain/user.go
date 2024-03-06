package domain

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id             uuid.UUID `json:"-" gorm:"primary key"`
	Username       string    `json:"username" binding:"required"`
	Email          string    `json:"email" binding:"required,email" gorm:"unique"`
	Password       string    `json:"password" binding:"required"`
	IsAdmin        bool      `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

//	ProfilePicture
