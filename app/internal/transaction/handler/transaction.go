package handler

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/shared"
	"bsnack/app/internal/transaction/dto"
	httphelper "bsnack/app/pkg/http"
	"bsnack/app/pkg/validation"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type TransactionHandlerImpl struct {
	TransactionUseCase interfaces.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase interfaces.TransactionUseCase, productUseCase interfaces.ProductUseCase) interfaces.TransactionHandler {
	return &TransactionHandlerImpl{
		TransactionUseCase: transactionUseCase,
	}
}

func (t *TransactionHandlerImpl) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, total, err := t.TransactionUseCase.GetTransactions(r.Context())
	if err != nil {
		httphelper.HandleError(w, err)
		return
	}

	transactionResponse := make([]dto.GetTransactionResponse, len(*transactions))
	for i, transaction := range *transactions {
		transactionResponse[i] = dto.GetTransactionResponse{
			ID:            transaction.ID,
			CustomerName:  transaction.CustomerName,
			ProductName:   transaction.ProductName,
			ProductFlavor: transaction.ProductFlavor,
			ProductSize:   transaction.ProductSize,
			Quantity:      transaction.Quantity,
			CreatedAt:     transaction.CreatedAt,
		}
	}

	httphelper.JSONResponse(w, http.StatusOK, "Transactions retrieved successfully", shared.PaginatedResponse[dto.GetTransactionResponse]{
		Data:  transactionResponse,
		Total: total,
	})
}

func (t *TransactionHandlerImpl) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	transactionId := chi.URLParam(r, "id")
	if transactionId == "" {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Transaction ID is required. Path: "+r.URL.Path, nil)
		return
	}

	transaction, err := t.TransactionUseCase.GetTransactionById(r.Context(), transactionId)
	if err != nil {
		httphelper.HandleError(w, err)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Transaction retrieved successfully", transaction)
}

func (t *TransactionHandlerImpl) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction dto.CreateTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid transaction data", nil)
		return
	}

	if err := validation.Validate.Struct(transaction); err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid transaction data", nil)
		return
	}

	createdTransaction, err := t.TransactionUseCase.CreateTransaction(r.Context(), &transaction)
	if err != nil {
		httphelper.HandleError(w, err)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Transaction created successfully", createdTransaction)
}
