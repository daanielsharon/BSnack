package dto

import (
	"time"
)

type CreateProductResponse struct {
	Name            string
	Type            string
	Flavor          string
	Size            string
	Price           float64
	ManufactureDate time.Time
	CreatedAt       time.Time
}

type GetProductResponse struct {
	Name   string
	Type   string
	Flavor string
	Size   string
	Price  float64
}
