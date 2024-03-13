package domain

import (
	"mime/multipart"
	"time"
)

type Products struct {
	Id           int       `json:"-"`
	MerchantId   int       `json:"-"`
	CategoryId   int       `json:"-"`
	Name         string    `json:"name"`
	Price        uint      `json:"price"`
	Description  string    `json:"description"`
	ProductPhoto string    `json:"product_photo"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type ProductParam struct {
	Id         int    `json:"-"`
	MerchantId int    `json:"-"`
	CategoryId int    `json:"-"`
	Name       string `json:"-"`
}

type ProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Price       uint   `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

type ProductUpdate struct {
	Name        string `json:"name"`
	Price       uint   `json:"price"`
	Description string `json:"description"`
}

type UploadProductPhoto struct {
	ProductPhoto *multipart.FileHeader `json:"product_photo"`
}
