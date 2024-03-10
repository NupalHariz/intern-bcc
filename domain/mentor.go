package domain

import "time"

type Mentors struct {
	Id            int       `json:"-"`
	Name          string    `json:"-"`
	CurrentJob    string    `json:"-"`
	Description   string    `json:"-"`
	MentorPicture string    `json:"-"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

type MentorRequest struct {
	Name          string `json:"name" binding:"required"`
	CurrentJob    string `json:"current_job" binding:"required"`
	Description   string `json:"description"`
	MentorPicture string `json:"-"`
}
