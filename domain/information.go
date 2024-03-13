package domain

import "time"

type Information struct {
	Id               int       `json:"-"`
	Title            string    `json:"title"`
	CategoryId       int       `json:"-"`
	Synopsis         string    `json:"synopsis"`
	Content          string    `json:"content"`
	InformationPhoto string    `json:"information_photo"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

type InformationRequest struct {
	Title    string `json:"title" binding:"required"`
	Category string `json:"category" binding:"required"`
	Synopsis string `json:"synopsis"`
	Content  string `json:"content"`
}

type InformationUpdate struct {
	Sysnopsis string `json:"synopsis"`
	Content   string `json:"content"`
}
