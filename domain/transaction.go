package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transactions struct {
	Id          uuid.UUID `json:"-" gorm:"primary key"`
	UserId      uuid.UUID `json:"-"`
	MentorId    int       `json:"-"`
	Price       uint64    `json:"-"`
	IsPayed     bool    `json:"-"`
	PaymentType string    `json:"-"`
	CreatedAt   time.Time `json:"-"`
	PayedAt     time.Time `json:"-" gorm:"autoCreateTime;autoUpdateTime"`
	Mentor      Mentors   `json:"-"`
	User        Users     `json:"-"`
}

type TransactionRequest struct {
	Price       uint64 `json:"price" binding:"required"`
	PaymentType string `json:"payment_type" binding:"required"`
}
