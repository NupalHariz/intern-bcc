package domain

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Products struct {
	Id           uuid.UUID `json:"-" gorm:"type:varchar(36);primary key"`
	MerchantId   uuid.UUID `json:"-" gorm:"type:varchar(36)"`
	CategoryId   int       `json:"-"`
	Name         string    `json:"name"`
	Price        uint      `json:"price"`
	Description  string    `json:"description"`
	ProductPhoto string    `json:"product_photo"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	LikeByUser   []Users   `json:"-" gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:product_id;references:id;joinReferences:user_id"`
	Merchant     Merchants `json:"merchant"`
}

type ProductParam struct {
	Id           uuid.UUID `json:"-"`
	MerchantId   int       `json:"-"`
	CategoryId   int       `json:"-" form:"categoryId"`
	Name         string    `json:"-" form:"name"`
	ProvinceId   int       `json:"-" form:"province"`
	UniversityId int       `json:"-" form:"university"`
	Page         int       `json:"-" form:"page" gorm:"-"`
	Offset       int       `json:"-" gorm:"-"`
}

type ProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Price       uint   `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

type ProductUpdate struct {
	Name         string `json:"name"`
	Price        uint   `json:"price"`
	Description  string `json:"description"`
	ProductPhoto string `json:"-"`
}

type UploadProductPhoto struct {
	ProductPhoto *multipart.FileHeader `json:"product_photo"`
}

type ProductResponses struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	MerchantName string    `json:"merchant_name"`
	University   string    `json:"university"`
	Price        uint      `json:"price"`
	ProductPhoto string    `json:"product_photo"`
}

type ProductResponse struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	MerchantName string    `json:"merchant_name"`
	University   string    `json:"university"`
	Faculty      string    `json:"faculty"`
	Province     string    `json:"province"`
	City         string    `json:"city"`
	Price        uint      `json:"price"`
	ProductPhoto string    `json:"product_photo"`
	WhatsApp     string    `json:"whatsapp"`
	Instagram    string    `json:"instagram"`
}
