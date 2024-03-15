package domain

type Province struct {
	Id        int         `json:"-"`
	Province  string      `json:"province"`
	Merchants []Merchants `json:"-" gorm:"foreignKey:province_id;references:id"`
}
