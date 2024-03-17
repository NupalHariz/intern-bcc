package database

import "intern-bcc/domain"

func Migrate() {
	DB.Migrator().AutoMigrate(
		&domain.Users{},
		&domain.Transactions{},
		&domain.Universities{},
		&domain.Province{},
		&domain.Information{},
		&domain.Merchants{},
		&domain.Products{},
		&domain.Mentors{},
		&domain.Experiences{},
		&domain.Categories{},
	)
}
