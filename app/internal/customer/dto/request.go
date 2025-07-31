package dto

type CreateCustomerRequest struct {
	Name   string `json:"name" validate:"required"`
	Points int    `json:"points" validate:"min=0"`
}

type CreatePointRedemptionRequest struct {
	ProductName   string `json:"product_name" validate:"required"`
	ProductSize   string `json:"product_size" validate:"required,product_size"`
	ProductFlavor string `json:"product_flavor" validate:"required,product_flavor"`
	Quantity      int    `json:"quantity" validate:"min=1"`
}
