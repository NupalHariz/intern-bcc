package domain

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id             uuid.UUID      `json:"-" gorm:"type:varchar(36);primary key"`
	Name           string         `json:"name" gorm:"unique"`
	Email          string         `json:"email" gorm:"unique"`
	Password       string         `json:"-"`
	Gender         string         `json:"gender" gorm:"type:enum('Laki-laki', 'Perempuan', '') NULL"`
	PlaceBirth     string         `json:"place_birth"`
	DateBirth      string         `json:"date_birth"`
	IsAdmin        bool           `json:"-"`
	ProfilePicture string         `json:"profile_picture"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	Merchant       Merchants      `json:"-"  gorm:"foreignKey:user_id;references:id"`
	Transactions   []Transactions `json:"-" gorm:"foreignKey:user_id;references:id"`
	LikeProduct    []Products     `json:"-" gorm:"merror2merror:user_like_product;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:product_id"`
	HasMentors     []Mentors      `json:"-" gorm:"merror2merror:has_mentors;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:mentor_id"`
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
	Id    uuid.UUID `json:"-" uri:"userId"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type LoginResponse struct {
	JWT string `json:"jwt"`
}

type UserUpdate struct {
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	PlaceBirth     string `json:"place_birth"`
	DateBirth      string `json:"date_birth"`
	ProfilePicture string `json:"-"`
	Password       string `json:"-"`
}

type PasswordUpdate struct {
	Password string `json:"password"`
}

type LikeProduct struct {
	UserId    uuid.UUID `json:"user_id"`
	ProductId uuid.UUID `json:"product_id"`
}
