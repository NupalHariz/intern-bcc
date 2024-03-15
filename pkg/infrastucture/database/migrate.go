package database

import "intern-bcc/domain"

func Migrate() {
	DB.AutoMigrate(
		&domain.Universities{},
		&domain.Province{},
		&domain.Categories{},
		&domain.Users{},
		&domain.Products{},
		&domain.Merchants{},
		&domain.Mentors{},
		&domain.Transactions{},
		&domain.Experiences{},
		&domain.Information{},
	)
}
