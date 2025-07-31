package dto

type CreateProductRequest struct {
	Name            string  `json:"name" validate:"required"`
	Type            string  `json:"type" validate:"required,product_type"`
	Flavor          string  `json:"flavor" validate:"required,product_flavor"`
	Size            string  `json:"size" validate:"required,product_size"`
	Price           float64 `json:"price" validate:"required,product_price"`
	Quantity        int     `json:"quantity" validate:"required,min=1"`
	ManufactureDate string  `json:"manufacture_date" validate:"required,YYYY-MM-DD_dateFormat"`
}
