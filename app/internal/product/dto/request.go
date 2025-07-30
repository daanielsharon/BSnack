package dto

import "time"

type CreateProductRequest struct {
	Name            string    `json:"name" validate:"required"`
	Type            string    `json:"type" validate:"required"`
	Flavor          string    `json:"flavor" validate:"required"`
	Size            string    `json:"size" validate:"required"`
	Price           float64   `json:"price" validate:"required"`
	Quantity        int       `json:"quantity" validate:"required"`
	ManufactureDate time.Time `json:"manufacture_date" validate:"required,YYYY-MM-DD_dateFormat"`
}
