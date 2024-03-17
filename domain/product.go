package domain

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Products struct {
	Id           uuid.UUID `json:"-"`
	MerchantId   uuid.UUID `json:"-"`
	CategoryId   int       `json:"-"`
	Name         string    `json:"name"`
	Price        uint      `json:"price"`
	Description  string    `json:"description"`
	ProductPhoto string    `json:"product_photo"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	LikeByUser   []Users   `json:"-" gorm:"merror2merror:user_like_product;foreignKey:id;joinForeignKey:product_id;references:id;joinReferences:user_id"`
}

type ProductParam struct {
	Id         uuid.UUID `json:"-"`
	MerchantId int       `json:"-"`
	CategoryId int       `json:"-"`
	Name       string    `json:"-"`
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
