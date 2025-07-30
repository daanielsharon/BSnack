package models

import (
	"time"
)

type Product struct {
	Name            string
	Type            string
	Flavor          string
	Size            string
	Price           float64
	Quantity        int
	ManufactureDate time.Time
	CreatedAt       time.Time
}
