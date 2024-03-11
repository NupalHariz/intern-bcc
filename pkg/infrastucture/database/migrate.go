package database

import "intern-bcc/domain"

func Migrate() {
	DB.AutoMigrate(
		&domain.Users{},
		&domain.Categories{},
		&domain.Products{},
		&domain.Merchants{},
		&domain.Mentors{},
		&domain.Transactions{},
	)
}
