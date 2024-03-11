package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type ITransactionRepository interface{
	CreateTransaction(newTransaction *domain.Transactions) error
}

type TransactionsRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionsRepository{db}
}

func (r *TransactionsRepository) CreateTransaction(newTransaction *domain.Transactions) error {
	tx := r.db.Begin()

	err := r.db.Create(newTransaction).Error
	if err != nil{
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
