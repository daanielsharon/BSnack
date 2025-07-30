package handler

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/transaction/dto"
	httphelper "bsnack/app/pkg/http"
	"bsnack/app/pkg/validation"
	"encoding/json"
	"net/http"
)

type TransactionHandlerImpl struct {
	TransactionUseCase interfaces.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase interfaces.TransactionUseCase) interfaces.TransactionHandler {
	return &TransactionHandlerImpl{
		TransactionUseCase: transactionUseCase,
	}
}

func (t *TransactionHandlerImpl) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := t.TransactionUseCase.GetTransactions(r.Context())
	if err != nil {
		httphelper.JSONResponse(w, http.StatusInternalServerError, "Failed to get transactions", nil)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Transactions retrieved successfully", transactions)
}

func (t *TransactionHandlerImpl) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	transactionId := r.URL.Query().Get("id")
	transaction, err := t.TransactionUseCase.GetTransactionById(r.Context(), transactionId)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusInternalServerError, "Failed to get transaction", nil)
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
		httphelper.JSONResponse(w, http.StatusInternalServerError, "Failed to create transaction", nil)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Transaction created successfully", createdTransaction)
}
