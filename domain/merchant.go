package domain

import (
	"time"

	"github.com/google/uuid"
)

type Merchants struct {
	Id          int       `json:"id"`
	UserId      uuid.UUID `json:"user_id" gorm:"type:varchar(36);unique"`
	StoreName   string    `json:"store_name"`
	University  string    `json:"university"`
	Faculty     string    `json:"faculty"`
	Province    string    `json:"province"`
	City        string    `json:"city"`
	PhoneNumber string    `json:"phone_number"`
	Instagram   string    `json:"instagram"`
	StorePhoto  string    `json:"-"`
	IsActive    bool      `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type MerchantRequest struct {
	StoreName   string `json:"store_name" `
	University  string `json:"university" binding:"required"`
	Faculty     string `json:"faculty" binding:"required"`
	Province    string `json:"province" binding:"required"`
	City        string `json:"city" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Instagram   string `json:"instagram"`
}

type MerchantParam struct {
	Id          int       `json:"id"`
	UserId      uuid.UUID `json:"user_id" gorm:"type:varchar(36);unique"`
	StoreName   string    `json:"store_name" `
	University  string    `json:"university" binding:"required"`
	Faculty     string    `json:"faculty" binding:"required"`
	Province    string    `json:"province" binding:"required"`
	City        string    `json:"city" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Instagram   string    `json:"instagram"`
}

type MerchantVerify struct {
	VerifyOtp string `json:"verify_otp" binding:"required,min=6,max=6"`
}
