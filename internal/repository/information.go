package repository

import (
	"fmt"
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IInformationRepository interface {
	GetArticles(articles *[]domain.Articles) error
	GetWebinarNCompetition(webinarNCompetition *[]domain.Information) error
	GetInformation(information *domain.Information, informationParam domain.InformationParam) error
	CreateInformation(newInformation *domain.Information) error
	UpdateInformation(information *domain.InformationUpdate, informationId int) error
}

type InformationRepository struct {
	db *gorm.DB
}

func NewInformationRepository(db *gorm.DB) IInformationRepository {
	return &InformationRepository{db}
}

func (r *InformationRepository) GetArticles(articles *[]domain.Articles) error {
	err := r.db.Model(domain.Information{}).Where("category_id = ?", 7).Order("created_at desc").Limit(15).Find(articles).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *InformationRepository) GetWebinarNCompetition(webinarNCompetition *[]domain.Information) error {
	err := r.db.Model(domain.Information{}).Where("category_id IN (?)", []int64{8, 9}).Limit(15).Order("created_at desc").Preload("Category").Find(webinarNCompetition).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *InformationRepository) GetInformation(information *domain.Information, informationParam domain.InformationParam) error {
	fmt.Println(informationParam)
	err := r.db.First(information, informationParam).Error
	if err != nil {
		return nil
	}

	return err
}

func (r *InformationRepository) CreateInformation(newInformation *domain.Information) error {
	tx := r.db.Begin()

	err := r.db.Create(newInformation).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *InformationRepository) UpdateInformation(information *domain.InformationUpdate, informationId int) error {
	err := r.db.Model(domain.Information{}).Where("id = ?", informationId).Updates(information).Error
	if err != nil {
		return err
	}

	return nil
}
