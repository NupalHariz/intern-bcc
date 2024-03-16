package domain

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Mentors struct {
	Id            uuid.UUID      `json:"-"`
	Name          string         `json:"name" gorm:"unique"`
	CurrentJob    string         `json:"current_job"`
	Description   string         `json:"description"`
	Price         uint64         `json:"price"`
	MentorPicture string         `json:"mentor_picture"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	Transactions  []Transactions `json:"-" gorm:"foreignKey:mentor_id;references:id"`
	Experiences   []Experiences  `json:"-" gorm:"foreignKey:mentor_id;references:id"`
	Users         []Users        `json:"-" gorm:"many2many:has_mentors;foreignKey:id;joinForeignKey:mentor_id;references:id;joinReferences:user_id"`
}

type MentorRequest struct {
	Name        string `json:"name" binding:"required"`
	CurrentJob  string `json:"current_job" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       uint64 `json:"price" binding:"required"`
}

type MentorUpdate struct {
	CurrentJob  string `json:"current_job" binding:"required"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
}

type UploadMentorPicture struct {
	MentorPicture *multipart.FileHeader `json:"mentor_picture"`
}

type HasMentor struct {
	UserId   uuid.UUID `json:"-"`
	MentorId int       `json:"-"`
}
