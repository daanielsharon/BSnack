package usecase

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/internal/transaction/dto"
	"bsnack/app/internal/validation"
	httphelper "bsnack/app/pkg/http"
	"math"
	"net/http"

	"context"

	"gorm.io/gorm"
)

type TransactionUseCaseImpl struct {
	transactionRepository interfaces.TransactionRepository
	customerUseCase       interfaces.CustomerUseCase
	productUseCase        interfaces.ProductUseCase
}

func NewTransactionUseCase(transactionRepository interfaces.TransactionRepository, customerUseCase interfaces.CustomerUseCase, productUseCase interfaces.ProductUseCase) interfaces.TransactionUseCase {
	return &TransactionUseCaseImpl{
		transactionRepository: transactionRepository,
		customerUseCase:       customerUseCase,
		productUseCase:        productUseCase,
	}
}

func (t *TransactionUseCaseImpl) GetTransactions(ctx context.Context) (*[]models.Transaction, error) {
	return t.transactionRepository.GetTransactions(ctx)
}

func (t *TransactionUseCaseImpl) GetTransactionById(ctx context.Context, id string) (*dto.GetTransactionResponse, error) {
	transaction, err := t.transactionRepository.GetTransactionById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, httphelper.NewAppError(http.StatusNotFound, "Transaction not found")
		}
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
	product, err := t.productUseCase.GetProductByName(ctx, transaction.ProductName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, httphelper.NewAppError(http.StatusNotFound, "Product not found")
		}
		return nil, err
	}

	if err := validation.ValidateProductExists(product); err != nil {
		return nil, err
	}

	if err := validation.ValidateProductStockEnough(product, transaction.Quantity); err != nil {
		return nil, err
	}

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

	pointsAdded := math.RoundToEven(product.Price/1000) * float64(transaction.Quantity)
	updatedCustomer, err := t.customerUseCase.AddCustomerPoint(ctx, transaction.CustomerName, int(pointsAdded))
	if err != nil {
		return nil, err
	}

	return &dto.CreateTransactionResponse{
		ID:            createdTransaction.ID,
		CustomerName:  updatedCustomer.Name,
		ProductName:   createdTransaction.ProductName,
		ProductFlavor: createdTransaction.ProductFlavor,
		ProductSize:   createdTransaction.ProductSize,
		Quantity:      createdTransaction.Quantity,
		CreatedAt:     createdTransaction.CreatedAt,
	}, nil
}
