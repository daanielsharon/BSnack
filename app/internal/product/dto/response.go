package dto

import (
	"time"
)

type CreateProductResponse struct {
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	Flavor          string    `json:"flavor"`
	Size            string    `json:"size"`
	Price           float64   `json:"price"`
	ManufactureDate string    `json:"manufacture_date"`
	CreatedAt       time.Time `json:"created_at"`
}

type GetProductResponse struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Flavor   string  `json:"flavor"`
	Size     string  `json:"size"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
