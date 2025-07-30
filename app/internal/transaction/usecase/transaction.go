package usecase

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/internal/transaction/dto"

	"context"
)

type TransactionUseCaseImpl struct {
	transactionRepository interfaces.TransactionRepository
}

func NewTransactionUseCase(transactionRepository interfaces.TransactionRepository) interfaces.TransactionUseCase {
	return &TransactionUseCaseImpl{
		transactionRepository: transactionRepository,
	}
}

func (t *TransactionUseCaseImpl) GetTransactions(ctx context.Context) (*[]dto.GetTransactionResponse, error) {
	transactions, err := t.transactionRepository.GetTransactions(ctx)
	if err != nil {
		return nil, err
	}

	var transactionResponse []dto.GetTransactionResponse
	for _, transaction := range *transactions {
		transactionResponse = append(transactionResponse, dto.GetTransactionResponse{
			ID:            transaction.ID,
			CustomerName:  transaction.CustomerName,
			ProductName:   transaction.ProductName,
			ProductFlavor: transaction.ProductFlavor,
			ProductSize:   transaction.ProductSize,
			Quantity:      transaction.Quantity,
			CreatedAt:     transaction.CreatedAt,
		})
	}
	return &transactionResponse, nil
}

func (t *TransactionUseCaseImpl) GetTransactionById(ctx context.Context, id string) (*dto.GetTransactionResponse, error) {
	transaction, err := t.transactionRepository.GetTransactionById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.GetTransactionResponse{
		ID:            transaction.ID,
		CustomerName:  transaction.CustomerName,
		ProductName:   transaction.ProductName,
		ProductFlavor: transaction.ProductFlavor,
		ProductSize:   transaction.ProductSize,
		Quantity:      transaction.Quantity,
		CreatedAt:     transaction.CreatedAt,
	}, nil
}

func (t *TransactionUseCaseImpl) CreateTransaction(ctx context.Context, transaction *dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {
	convertedTransaction := &models.Transaction{
		CustomerName:  transaction.CustomerName,
		ProductName:   transaction.ProductName,
		ProductFlavor: transaction.ProductFlavor,
		ProductSize:   transaction.ProductSize,
		Quantity:      transaction.Quantity,
	}

	createdTransaction, err := t.transactionRepository.CreateTransaction(ctx, convertedTransaction)
	if err != nil {
		return nil, err
	}

	return &dto.CreateTransactionResponse{
		ID:            createdTransaction.ID,
		CustomerName:  createdTransaction.CustomerName,
		ProductName:   createdTransaction.ProductName,
		ProductFlavor: createdTransaction.ProductFlavor,
		ProductSize:   createdTransaction.ProductSize,
		Quantity:      createdTransaction.Quantity,
		CreatedAt:     createdTransaction.CreatedAt,
	}, nil
}
