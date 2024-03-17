package domain

import "time"

type Experiences struct {
	MentorId   int       `json:"-"`
	Experience string    `json:"-"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type ExperienceRequest struct {
	Experience string `json:"experience" binding:"required"`
}
