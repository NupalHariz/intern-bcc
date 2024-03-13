package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IProductRepository interface{
	CreateProduct(newProduct *domain.Products) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{db}
}

func(r *ProductRepository) CreateProduct(newProduct *domain.Products) error {
	tx := r.db.Begin()

	err := r.db.Create(newProduct).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}