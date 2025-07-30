package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTransactionResponse struct {
	ID            uuid.UUID `json:"id"`
	CustomerName  string    `json:"customer_name"`
	ProductName   string    `json:"product_name"`
	ProductSize   string    `json:"product_size"`
	ProductFlavor string    `json:"product_flavor"`
	Quantity      int       `json:"quantity"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetTransactionResponse struct {
	ID            uuid.UUID `json:"id"`
	CustomerName  string    `json:"customer_name"`
	ProductName   string    `json:"product_name"`
	ProductSize   string    `json:"product_size"`
	ProductFlavor string    `json:"product_flavor"`
	Quantity      int       `json:"quantity"`
	CreatedAt     time.Time `json:"created_at"`
}
