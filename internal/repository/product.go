package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type IProductRepository interface {
	GetProduct(product *domain.Products, productParam *domain.ProductParam) error
	CreateProduct(newProduct *domain.Products) error
	UpdateProduct(product *domain.Products) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) GetProduct(product *domain.Products, productParam *domain.ProductParam) error {
	err := r.db.First(product, productParam).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) CreateProduct(newProduct *domain.Products) error {
	tx := r.db.Begin()

	err := r.db.Create(newProduct).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *ProductRepository) UpdateProduct(product *domain.Products) error {
	tx := r.db.Begin()

	err := r.db.Where("id = ?", product.Id).Updates(product).Error
	if err != nil {
		tx.Rollback()
		return nil
	}

	tx.Commit()
	return nil
}
