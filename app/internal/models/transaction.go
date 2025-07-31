package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid;not null;default:uuid_generate_v4()"`
	CustomerName  string    `gorm:"type:varchar(255);not null"`
	ProductName   string    `gorm:"type:varchar(255);not null"`
	ProductSize   string    `gorm:"type:varchar(255);not null"`
	ProductFlavor string    `gorm:"type:varchar(255);not null"`
	Quantity      int       `gorm:"type:int;not null"`
	CreatedAt     time.Time `gorm:"type:timestamp;not null"`
}
