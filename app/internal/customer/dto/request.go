package dto

type CreateCustomerRequest struct {
	Name   string `json:"name" validate:"required"`
	Points int    `json:"points" validate:"min=0"`
}

type CreatePointRedemptionRequest struct {
	CustomerName  string `json:"customer_name" validate:"required"`
	ProductName   string `json:"product_name" validate:"required,product_name"`
	ProductSize   string `json:"product_size" validate:"required,product_size"`
	ProductType   string `json:"product_type" validate:"required,product_type"`
	ProductFlavor string `json:"product_flavor" validate:"required,product_flavor"`
	Quantity      int    `json:"quantity" validate:"min=1"`
}
