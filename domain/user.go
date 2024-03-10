package domain

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id             uuid.UUID `json:"-" gorm:"type:varchar(36);primary key"`
	Username       string    `json:"username"`
	Email          string    `json:"email" gorm:"unique"`
	Password       string    `json:"password"`
	IsAdmin        bool      `json:"-"`
	ProfilePicture string    `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	Merchant       Merchants `json:"-"  gorm:"foreignKey:user_id;references:id"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"is_admin"`
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
