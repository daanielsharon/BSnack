package models

import (
	"time"
)

type Product struct {
	Name            string    `gorm:"primarykey;type:varchar(255);not null"`
	Type            string    `gorm:"type:varchar(255);not null"`
	Flavor          string    `gorm:"type:varchar(255);not null"`
	Size            string    `gorm:"type:varchar(255);not null"`
	Price           float64   `gorm:"type:decimal(10,2);not null"`
	Quantity        int       `gorm:"type:int;not null"`
	ManufactureDate time.Time `gorm:"type:timestamp;not null"`
	CreatedAt       time.Time `gorm:"type:timestamp;not null"`
}
