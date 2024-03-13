package domain

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id             uuid.UUID      `json:"-" gorm:"type:varchar(36);primary key"`
	Name           string         `json:"name"`
	Email          string         `json:"email" gorm:"unique"`
	Password       string         `json:"-"`
	Gender         string         `json:"gender" gorm:"type:enum('Laki-laki', 'Perempuan', '')"`
	PlaceBirth     string         `json:"place_birth"`
	DateBirth      string         `json:"date_birth"`
	IsAdmin        bool           `json:"-"`
	ProfilePicture string         `json:"profile_picture"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	Merchant       Merchants      `json:"-"  gorm:"foreignKey:user_id;references:id"`
	Transactions   []Transactions `json:"-" gorm:"foreignKey:user_id;references:id"`
}

type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserParam struct {
	Id       uuid.UUID `json:"-"`
	Email    string    `json:"-"`
}

type LoginResponse struct {
	JWT string `json:"jwt"`
}

type UserUpdate struct {
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	PlaceBirth string `json:"place_birth"`
	DateBirth  string `json:"date_birth"`
}

type UploadUserPhoto struct {
	ProfilePicture *multipart.FileHeader `json:"profile_picture"`
}
