package domain

import (
	"time"
)

type Products struct {
	Id           int       `json:"-"`
	MerchantId   int       `json:"-"`
	CategoryId   int       `json:"-"`
	Name         string    `json:"-"`
	Price        uint      `json:"-"`
	Description  string    `json:"-"`
	ProductPhoto string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type ProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Price       uint   `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

//Product Photo
