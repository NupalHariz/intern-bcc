package domain

import (
	"time"
)

type Products struct {
	Id           int       `json:"id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Price        uint      `json:"price" binding:"required"`
	ProductPhoto string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

//Product Photo
