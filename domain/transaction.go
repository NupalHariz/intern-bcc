package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transactions struct {
	Id          uuid.UUID `json:"-" gorm:"primary key"`
	UserId      uuid.UUID `json:"-"`
	MentorId    uuid.UUID `json:"-"`
	Price       uint64    `json:"-"`
	IsPayed     bool      `json:"-"`
	PaymentType string    `json:"-"`
	CreatedAt   time.Time `json:"-"`
	PayedAt     time.Time `json:"-" gorm:"type:datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP"`
	Mentor      Mentors   `json:"-"`
	User        Users     `json:"-"`
}

type TransactionRequest struct {
	Price       uint64 `json:"price" binding:"required"`
	PaymentType string `json:"payment_type" binding:"required"`
}

type TransactionResponse struct {
	OrderId     string `json:"order_id"`
	PaymentType string `json:"payment_type"`
	VaNumber    string `json:"va_number"`
	BillerCode  string `json:"biller_code"`
	BillKey     string `json:"bill_key"`
	URL         string `json:"url"`
}
