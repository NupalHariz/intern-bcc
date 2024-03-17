package domain

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Merchants struct {
	Id            uuid.UUID    `json:"id" gorm:"type:varchar(36);unique"`
	UserId        uuid.UUID    `json:"user_id" gorm:"type:varchar(36);unique"`
	MerchantName  string       `json:"store_name"`
	UniversityId  int          `json:"-"`
	Faculty       string       `json:"faculty"`
	ProvinceId    int          `json:"-"`
	City          string       `json:"city"`
	PhoneNumber   string       `json:"phone_number"`
	Instagram     string       `json:"instagram"`
	MerchantPhoto string       `json:"merchant_photo"`
	IsActive      bool         `json:"-"`
	CreatedAt     time.Time    `json:"-"`
	UpdatedAt     time.Time    `json:"-"`
	Products      []Products   `json:"-" gorm:"foreignKey:merchant_id;references:id"`
	University    Universities `json:"university"`
}

type MerchantRequest struct {
	MerchantName string `json:"store_name" `
	University   string `json:"university" binding:"required"`
	Faculty      string `json:"faculty" binding:"required"`
	Province     string `json:"province" binding:"required"`
	City         string `json:"city" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	Instagram    string `json:"instagram"`
}

type MerchantParam struct {
	Id           uuid.UUID `json:"id"`
	UserId       uuid.UUID `json:"user_id"`
	MerchantName string    `json:"store_name"`
}

type MerchantVerify struct {
	VerifyOtp string `json:"verify_otp" binding:"required,min=6,max=6"`
}

type UpdateMerchant struct {
	MerchantName  string `json:"store_name"`
	UniversityId  int    `json:"-"`
	Faculty       string `json:"-"`
	ProvinceId    int    `json:"-"`
	City          string `json:"-"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	Instagram     string `json:"instagram"`
	MerchantPhoto string `json:"-"`
	IsActive      bool   `json:"-"`
}

type UploadMerchantPhoto struct {
	MerchantPhoto *multipart.FileHeader `json:"merchant_photo"`
}
