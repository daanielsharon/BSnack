package repository

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"

	"context"

	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) interfaces.TransactionRepository {
	return &TransactionRepositoryImpl{DB: db}
}

func (t *TransactionRepositoryImpl) GetTransactions(ctx context.Context) (*[]models.Transaction, error) {
	var transactions []models.Transaction
	err := t.DB.WithContext(ctx).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return &transactions, nil
}

func (t *TransactionRepositoryImpl) GetTransactionById(ctx context.Context, id string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := t.DB.WithContext(ctx).Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (t *TransactionRepositoryImpl) CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	return transaction, t.DB.WithContext(ctx).Model(&models.Transaction{}).Create(transaction).Error
}
