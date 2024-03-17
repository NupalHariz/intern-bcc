package repository

import (
	"intern-bcc/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IProductRepository interface {
	GetProduct(product *domain.Products, productParam domain.ProductParam) error
	CreateProduct(newProduct *domain.Products) error
	UpdateProduct(product *domain.ProductUpdate, productId uuid.UUID) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) GetProduct(product *domain.Products, productParam domain.ProductParam) error {
	err := r.db.First(product, productParam).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) CreateProduct(newProduct *domain.Products) error {
	err := r.db.Create(newProduct).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) UpdateProduct(product *domain.ProductUpdate, productId uuid.UUID) error {
	err := r.db.Where("id = ?", productId).Updates(product).Error
	if err != nil {
		return nil
	}

	return nil
}
