package infrastucture

import (
	"fmt"
	"intern-bcc/domain"
	"log"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	var totalCategory int64

	if err := db.Debug().Model(domain.Categories{}).Count(&totalCategory).Error; err != nil {
		log.Fatal("error occured when counting category: ", err)
	}

	if totalCategory == 0 {
		if err := generateCategory(db); err != nil {
			log.Fatal("error occured when creating book: ", err)
		}
	}

	var totalAdmin int64
	if err := db.Model(domain.Users{}).Where("is_admin = ?", true).Count(&totalAdmin).Error; err != nil {
		log.Fatal("error occured when counting admin: ", err)
	}

	fmt.Println(totalAdmin)
	if totalAdmin == 0 {
		if err := generateAdmin(db); err != nil {
			log.Fatal("error occured when creating admin: ", err)
		}
	}

	var totalMentor int64
	if err := db.Model(domain.Mentors{}).Count(&totalMentor).Error; err != nil {
		log.Fatal("error occured when counting mentor: ", err)
	}

	if totalMentor == 0 {
		if err := generateMentor(db); err != nil {
			log.Fatal("error occured when creating mentor: ", err)
		}
	}

}

func generateMentor(db *gorm.DB) error {
	var mentors []*domain.Mentors

	for i := 1; i <= 5; i++ {
		mentor := &domain.Mentors{
			Name:        faker.Name(),
			CurrentJob:  faker.Word(),
			Description: faker.Sentence(),
		}

		mentors = append(mentors, mentor)
	}

	err := db.CreateInBatches(&mentors, 5).Error
	if err != nil {
		log.Fatal("error occured when creating mentor: ", err)
	}

	return nil
}

func generateAdmin(db *gorm.DB) error {
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("rahasiaadmin"), 10)
	if err != nil {
		log.Fatal("error occured when hashing admin password: ", err)
	}

	admin := domain.Users{
		Id:       uuid.New(),
		Name:     "Admin",
		Email:    "Admin@gmail.com",
		Password: string(adminPassword),
		IsAdmin:  true,
	}

	err = db.Model(domain.Users{}).Create(&admin).Error
	if err != nil {
		log.Fatal("error occured when creating admin: ", err)
	}

	return nil
}

func generateCategory(db *gorm.DB) error {
	var categories []*domain.Categories

	categories = append(categories,
		&domain.Categories{
			Category: "Makanan & Minuman",
		},
		&domain.Categories{
			Category: "Kosmetik",
		},
		&domain.Categories{
			Category: "Fashion",
		},
		&domain.Categories{
			Category: "Aksesoris",
		},
		&domain.Categories{
			Category: "Kerajinan",
		},
		&domain.Categories{
			Category: "Jasa",
		},
		&domain.Categories{
			Category: "Artikel",
		},
		&domain.Categories{
			Category: "Webinar",
		},
		&domain.Categories{
			Category: "Lomba",
		},
	)

	if err := db.CreateInBatches(&categories, 9).Error; err != nil {
		return err
	}

	return nil
}
