package dto

type CreateTransactionRequest struct {
	CustomerName  string `json:"customer_name" validate:"required"`
	ProductName   string `json:"product_name" validate:"required"`
	ProductSize   string `json:"product_size" validate:"required"`
	ProductFlavor string `json:"product_flavor" validate:"required"`
	Quantity      int    `json:"quantity" validate:"required"`
}
