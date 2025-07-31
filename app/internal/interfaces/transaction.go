package interfaces

import (
	"bsnack/app/internal/models"
	"bsnack/app/internal/transaction/dto"
	"context"
	"net/http"
)

type TransactionHandler interface {
	GetTransactions(w http.ResponseWriter, r *http.Request)
	GetTransactionById(w http.ResponseWriter, r *http.Request)
	CreateTransaction(w http.ResponseWriter, r *http.Request)
}

type TransactionUseCase interface {
	GetTransactions(ctx context.Context) (*[]models.Transaction, int64, error)
	GetTransactionById(ctx context.Context, id string) (*dto.GetTransactionResponse, error)
	CreateTransaction(ctx context.Context, transaction *dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error)
}

type TransactionRepository interface {
	GetTransactions(ctx context.Context) (*[]models.Transaction, int64, error)
	GetTransactionById(ctx context.Context, id string) (*models.Transaction, error)
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
}
