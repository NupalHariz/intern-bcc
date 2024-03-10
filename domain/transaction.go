package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transactions struct {
	Id        uuid.UUID `json:"-"`
	UserId    uuid.UUID `json:"-"`
	MentorId  int       `json:"-"`
	Price     uint64    `json:"-"`
	Status    bool      `json:"-"`
	CreatedAt time.Time `json:"-"`
	PayedAt   time.Time `json:"-"`
}
