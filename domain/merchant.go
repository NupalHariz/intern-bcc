package domain

import "time"

type Merchants struct {
	Id           int	`json:"id" binding:"required"`
	University   string	`json:"university" binding:"required"`
	Faculty      string	`json:"faculty" binding:"required"`
	Province     string	`json:"province" binding:"required"`
	City         string	`json:"city" binding:"required"`
	PhoneNumber  string	`json:"phone_number" binding:"required"`
	Instagram    string	`json:"instagram"`
	IsActive     bool	`json:"-"`
	CreatedAt    time.Time	`json:"-"`
	UpdatedAt    time.Time	`json:"-"`
}

//Store Photo
