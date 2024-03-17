package domain

type Universities struct {
	Id         int         `json:"-"`
	University string      `json:"university"`
	Merchants  []Merchants `json:"-" gorm:"foreignKey:university_id;references:id"`
}

type UniversityRequest struct {
	University string `json:"university" binding:"required"`
}