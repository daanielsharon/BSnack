package interfaces

import (
	"bsnack/app/internal/customer/dto"
	"bsnack/app/internal/models"
	"context"
	"net/http"
)

type CustomerHandler interface {
	GetCustomers(w http.ResponseWriter, r *http.Request)
	GetCustomerByName(w http.ResponseWriter, r *http.Request)
	CreateCustomer(w http.ResponseWriter, r *http.Request)
}

type CustomerUseCase interface {
	GetCustomers(ctx context.Context) (*[]dto.GetCustomerResponse, error)
	GetCustomerByName(ctx context.Context, name string) (*dto.GetCustomerResponse, error)
	CreateCustomer(ctx context.Context, customer *dto.CreateCustomerRequest) (*dto.CreateCustomerResponse, error)
}

type CustomerRepository interface {
	GetCustomers(ctx context.Context) (*[]models.Customer, error)
	GetCustomerByName(ctx context.Context, name string) (*models.Customer, error)
	CreateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error)
}
