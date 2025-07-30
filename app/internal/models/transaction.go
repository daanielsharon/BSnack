package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID
	CustomerName  string
	ProductName   string
	ProductSize   string
	ProductFlavor string
	Quantity      int
	CreatedAt     time.Time
}
