package domain

import (
	"time"

	"github.com/google/uuid"
)

type Experiences struct {
	MentorId   uuid.UUID `json:"-"`
	Experience string    `json:"-"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type ExperienceRequest struct {
	Experience string `json:"experience" binding:"required"`
}
