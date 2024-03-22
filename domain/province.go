package domain

type Province struct {
	Id        int         `json:"-"`
	Province  string      `json:"province" gorm:"unique" binding:"required"`
	Merchants []Merchants `json:"-" gorm:"foreignKey:province_id;references:id"`
}
