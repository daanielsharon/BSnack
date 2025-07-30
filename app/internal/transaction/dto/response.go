package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTransactionResponse struct {
	ID            uuid.UUID
	CustomerName  string
	ProductName   string
	ProductSize   string
	ProductFlavor string
	Quantity      int
	CreatedAt     time.Time
}

type GetTransactionResponse struct {
	ID            uuid.UUID
	CustomerName  string
	ProductName   string
	ProductSize   string
	ProductFlavor string
	Quantity      int
	CreatedAt     time.Time
}
