package infrastucture

import (
	"intern-bcc/domain"
	"log"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	var totalCategory int64

	if err := db.Model(domain.Categories{}).Count(&totalCategory).Error; err != nil {
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

	var totalProvince int64
	if err := db.Model(domain.Province{}).Count(&totalProvince).Error; err != nil {
		log.Fatal("error occured when counting province: ", err)
	}

	if totalProvince == 0 {
		if err := generateProvince(db); err != nil {
			log.Fatal("error occured when creating province: ", err)
		}
	}

	var totalUniversity int64
	if err := db.Model(domain.Universities{}).Count(&totalUniversity).Error; err != nil {
		log.Fatal("error occured when counting university: ", err)
	}

	if totalUniversity == 0 {
		if err := generateUniversity(db); err != nil {
			log.Fatal("error occured when creating university: ", err)
		}
	}

	var totalInformation int64
	if err := db.Model(domain.Information{}).Count(&totalInformation).Error; err != nil {
		log.Fatal("error occured when counting Information: ", err)
	}

	if totalInformation == 0 {
		if err := generateInformation(db); err != nil {
			log.Fatal("error occured when creating informatin: ", err)
		}
	}
}

func generateInformation(db *gorm.DB) error {
	var informations []*domain.Information
	for i := 1; i < 15; i++ {
		article := domain.Information{
			Title:      faker.Name(),
			CategoryId: 7,
			Synopsis:   faker.Sentence(),
			Content:    faker.Paragraph(),
		}
		informations = append(informations, &article)

		webinar := domain.Information{
			Title:      faker.Name(),
			CategoryId: 8,
		}
		informations = append(informations, &webinar)

		lomba := domain.Information{
			Title:      faker.Name(),
			CategoryId: 9,
		}
		informations = append(informations, &lomba)
	}

	if err := db.CreateInBatches(&informations, 45).Error; err != nil {
		return err
	}

	return nil
}

func generateProvince(db *gorm.DB) error {
	var provinces []*domain.Province

	provinces = append(provinces,
		&domain.Province{
			Province: "Jawa Timur",
		},
		&domain.Province{
			Province: "Jawa Tengah",
		},
		&domain.Province{
			Province: "Jawa Barat",
		},
		&domain.Province{
			Province: "DKI Jakarta",
		},
		&domain.Province{
			Province: "DIY Yogyakarta",
		},
		&domain.Province{
			Province: "Sumatera Utara",
		},
		&domain.Province{
			Province: "Sumatera Barat",
		},
		&domain.Province{
			Province: "Kalimantan Utara",
		},
		&domain.Province{
			Province: "Kalimantan Barat",
		},
		&domain.Province{
			Province: "Bali",
		},
	)

	if err := db.CreateInBatches(&provinces, 10).Error; err != nil {
		return err
	}

	return nil
}

func generateUniversity(db *gorm.DB) error {
	var universities []*domain.Universities

	universities = append(universities,
		&domain.Universities{
			University: "Universitas Brawijaya",
		},
		&domain.Universities{
			University: "Universitas Airlangga",
		},
		&domain.Universities{
			University: "Universitas Negeri Semarang",
		},
		&domain.Universities{
			University: "Universitas Negeri Surabaya",
		},
		&domain.Universities{
			University: "Universitas Negeri Malang",
		},
		&domain.Universities{
			University: "Institut Teknologi Surabaya",
		},
		&domain.Universities{
			University: "Universitas Muhammadiyah Malang",
		},
		&domain.Universities{
			University: "Universitas Jember",
		},
		&domain.Universities{
			University: "Universitas Surabaya",
		},
		&domain.Universities{
			University: "Universitas Kristen Petra",
		},
	)

	if err := db.CreateInBatches(&universities, 10).Error; err != nil {
		return err
	}

	return nil
}

func generateMentor(db *gorm.DB) error {
	var mentors []*domain.Mentors

	for i := 1; i <= 5; i++ {
		mentor := &domain.Mentors{
			Id:          uuid.New(),
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
