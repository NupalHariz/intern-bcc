package repository

import (
	"intern-bcc/domain"

	"gorm.io/gorm"
)

type ITransactionRepository interface{
	GetTransaction(transaction *domain.Transactions) error
	CreateTransaction(newTransaction *domain.Transactions) error
	UpdateTransaction(transaction *domain.Transactions) error
}

type TransactionsRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionsRepository{db}
}

func (r *TransactionsRepository) GetTransaction(transaction *domain.Transactions) error {
	err := r.db.First(transaction, transaction).Error
	if err != nil {
		return err
	}

	return nil
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

func (r *TransactionsRepository) UpdateTransaction(transaction *domain.Transactions) error {
	tx := r.db.Begin()

	err := r.db.Where("id = ?", transaction.Id).Updates(domain.Transactions{
		PayedAt: transaction.PayedAt,
		IsPayed: transaction.IsPayed,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
