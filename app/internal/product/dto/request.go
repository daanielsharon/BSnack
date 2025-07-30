package dto

import "time"

type CreateProductRequest struct {
	Name            string    `json:"name" validate:"required"`
	Type            string    `json:"type" validate:"required,type"`
	Flavor          string    `json:"flavor" validate:"required,flavor"`
	Size            string    `json:"size" validate:"required,size"`
	Price           float64   `json:"price" validate:"required,price"`
	Quantity        int       `json:"quantity" validate:"required"`
	ManufactureDate time.Time `json:"manufacture_date" validate:"required,YYYY-MM-DD_dateFormat"`
}
