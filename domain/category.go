package domain

type Categories struct {
	Id          int           `json:"-"`
	Category    string        `json:"category" gorm:"unique" binding:"required"`
	Information []Information `json:"-" gorm:"foreignKey:category_id;references:id"`
	Product     []Products    `json:"-" gorm:"foreignKey:category_id;references:id"`
}
