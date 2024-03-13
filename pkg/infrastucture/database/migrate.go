package database

import "intern-bcc/domain"

func Migrate() {
	DB.AutoMigrate(
		&domain.Users{},
		&domain.Products{},
		&domain.Merchants{},
		&domain.Mentors{},
		&domain.Transactions{},
		&domain.Experiences{},
		&domain.Information{},
		&domain.Categories{},
	)
}
