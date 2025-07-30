package dto

type CreateCustomerRequest struct {
	Name   string `json:"name" validate:"required"`
	Points int    `json:"points" validate:"min=0"`
}
