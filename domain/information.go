package domain

import "time"

type Information struct {
	Id               int       `json:"-"`
	Title            string    `json:"-"`
	CategoryId       int       `json:"-"`
	Synopsis         string    `json:"-"`
	Content          string    `json:"-"`
	InformationPhoto string    `json:"-"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

type InformationRequest struct {
	Title    string `json:"title" binding:"required"`
	Category string `json:"category" binding:"required"`
	Synopsis string `json:"synopsis"`
	Content  string `json:"content"`
}
