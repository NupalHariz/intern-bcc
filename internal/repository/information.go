package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"intern-bcc/domain"

	"intern-bcc/pkg/redis"

	"gorm.io/gorm"
)

type IInformationRepository interface {
	GetArticles(ctx context.Context, articles *[]domain.Articles) error
	GetWebinarNCompetition(ctx context.Context, webinarNCompetition *[]domain.Information) error
	GetInformation(information *domain.Information, informationParam domain.InformationParam) error
	CreateInformation(newInformation *domain.Information) error
	UpdateInformation(information *domain.InformationUpdate, informationId int) error
}

type InformationRepository struct {
	db    *gorm.DB
	redis redis.IRedis
}

func NewInformationRepository(db *gorm.DB, redis redis.IRedis) IInformationRepository {
	return &InformationRepository{db, redis}
}

func (r *InformationRepository) GetArticles(ctx context.Context, articles *[]domain.Articles) error {
	key := fmt.Sprintf(KeySetInformationNmentor, "Articles")
	stringArticles, err := r.redis.GetRedis(ctx, key)
	if err != nil {
		err = r.db.Model(domain.Information{}).Where("category_id = ?", 7).Order("created_at desc").Limit(15).Find(articles).Error
		if err != nil {
			return err
		}

		byteArticles, err := json.Marshal(articles)
		if err != nil {
			return nil
		}

		err = r.redis.SetRedis(ctx, key, string(byteArticles), 5*time.Minute)
		if err != nil {
			return err
		}

		return nil
	}

	err = json.Unmarshal([]byte(stringArticles), articles)
	if err != nil {
		return err
	}
	return nil
}

func (r *InformationRepository) GetWebinarNCompetition(ctx context.Context, webinarNCompetition *[]domain.Information) error {
	key := fmt.Sprintf(KeySetInformationNmentor, "WebinarNCompetition")
	stringWebinerNCompetition, err := r.redis.GetRedis(ctx, key)
	if err != nil {
		err := r.db.Model(domain.Information{}).Where("category_id IN (?)", []int64{8, 9}).Limit(15).Order("created_at desc").Preload("Category").Find(webinarNCompetition).Error
		if err != nil {
			return err
		}

		byteWebinarNCompetition, err := json.Marshal(webinarNCompetition)
		if err != nil {
			return nil
		}

		err = r.redis.SetRedis(ctx, key, string(byteWebinarNCompetition), 5*time.Minute)
		if err != nil {
			return err
		}

		return nil
	}

	err = json.Unmarshal([]byte(stringWebinerNCompetition), webinarNCompetition)
	if err != nil {
		return nil
	}
	return nil
}

func (r *InformationRepository) GetInformation(information *domain.Information, informationParam domain.InformationParam) error {
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
