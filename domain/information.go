package domain

import "time"

type Information struct {
	Id               int        `json:"-"`
	Title            string     `json:"title"`
	CategoryId       int        `json:"-"`
	Synopsis         string     `json:"synopsis"`
	Content          string     `json:"content"`
	InformationPhoto string     `json:"information_photo"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	Category         Categories `json:"-"`
}

type InformationRequest struct {
	Title    string `json:"title" binding:"required"`
	Category string `json:"category" binding:"required"`
	Synopsis string `json:"synopsis"`
	Content  string `json:"content"`
}

type InformationUpdate struct {
	Sysnopsis        string `json:"synopsis"`
	Content          string `json:"content"`
	InformationPhoto string `json:"-"`
}

type InformationParam struct {
	Id int `json:"id"`
}

type Articles struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	Synopsis         string `json:"synopsis"`
	InformationPhoto string `json:"information_photo"`
}

type Article struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	InformationPhoto string `json:"information_photo"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type WebinarNCompetition struct {
	Id               int    `json:"id"`
	Title            string `json:"name"`
	Category         string `json:"category"`
	InformationPhoto string `json:"information_photo"`
}

type InformationResponses struct {
	Articles            []Articles            `json:"articles"`
	WebinarNCompetition []WebinarNCompetition `json:"webinar_and_competition"`
}
