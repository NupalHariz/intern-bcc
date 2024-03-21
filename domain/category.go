package domain

type Categories struct {
	Id          int           `json:"-"`
	Category    string        `json:"category"`
	Information []Information `json:"-" gorm:"foreignKey:category_id;references:id"`
	Product     []Products    `json:"-" gorm:"foreignKey:category_id;references:id"`
}

type CategoryRequest struct {
	Category string `json:"category" binding:"required"`
}
