package usecase

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/internal/transaction/dto"
	"bsnack/app/internal/validation"
	"bsnack/app/pkg/cache"
	httphelper "bsnack/app/pkg/http"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type TransactionUseCaseImpl struct {
	transactionRepository interfaces.TransactionRepository
	customerUseCase       interfaces.CustomerUseCase
	productUseCase        interfaces.ProductUseCase
	redisClient           *redis.Client
}

func NewTransactionUseCase(transactionRepository interfaces.TransactionRepository, customerUseCase interfaces.CustomerUseCase, productUseCase interfaces.ProductUseCase, redisClient *redis.Client) interfaces.TransactionUseCase {
	return &TransactionUseCaseImpl{
		transactionRepository: transactionRepository,
		customerUseCase:       customerUseCase,
		productUseCase:        productUseCase,
		redisClient:           redisClient,
	}
}

func (t *TransactionUseCaseImpl) GetTransactions(ctx context.Context) (*[]models.Transaction, int64, error) {
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
	product, err := t.productUseCase.GetProductByName(ctx, strings.ToLower(transaction.ProductName))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, httphelper.NewAppError(http.StatusNotFound, "Product not found")
		}
		return nil, err
	}

	if err := validation.ValidateSameProduct(product, &models.Product{
		Name:   strings.ToLower(transaction.ProductName),
		Flavor: strings.ToLower(transaction.ProductFlavor),
		Size:   strings.ToLower(transaction.ProductSize),
	}); err != nil {
		return nil, err
	}

	if err := validation.ValidateProductExists(product); err != nil {
		return nil, err
	}

	if err := validation.ValidateProductStockEnough(product, transaction.Quantity); err != nil {
		return nil, err
	}

	convertedTransaction := &models.Transaction{
		CustomerName:  strings.ToLower(transaction.CustomerName),
		ProductName:   strings.ToLower(transaction.ProductName),
		ProductFlavor: strings.ToLower(transaction.ProductFlavor),
		ProductSize:   strings.ToLower(transaction.ProductSize),
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

	if err := t.productUseCase.DeductProductStock(ctx, strings.ToLower(transaction.ProductName), transaction.Quantity); err != nil {
		return nil, err
	}

	err = cache.DeleteRedisKeysByPattern(ctx, t.redisClient, "customers:*")
	if err != nil {
		log.Printf("[WARN] Failed to delete cache for pattern %s in create transaction handler: %v", "customers:*", err)
	}

	err = cache.DeleteRedisKeysByPattern(ctx, t.redisClient, fmt.Sprintf("products:*:date=%s", product.GetDBManufactureDateInCorrectFormat()))
	if err != nil {
		log.Printf("[WARN] Failed to delete cache for pattern %s in create transaction handler: %v", fmt.Sprintf("products:*:date=%s", product.ManufactureDate), err)
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
