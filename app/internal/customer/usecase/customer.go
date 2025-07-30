package usecase

import (
	"bsnack/app/internal/customer/dto"
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"

	"context"
)

type CustomerUseCaseImpl struct {
	customerRepository interfaces.CustomerRepository
}

func NewCustomerUseCase(customerRepository interfaces.CustomerRepository) interfaces.CustomerUseCase {
	return &CustomerUseCaseImpl{
		customerRepository: customerRepository,
	}
}

func (c *CustomerUseCaseImpl) GetCustomers(ctx context.Context) (*[]dto.GetCustomerResponse, error) {
	customers, err := c.customerRepository.GetCustomers(ctx)
	if err != nil {
		return nil, err
	}

	var customerResponse []dto.GetCustomerResponse
	for _, customer := range *customers {
		customerResponse = append(customerResponse, dto.GetCustomerResponse{
			Name:   customer.Name,
			Points: customer.Points,
		})
	}
	return &customerResponse, nil
}

func (c *CustomerUseCaseImpl) GetCustomerByName(ctx context.Context, name string) (*dto.GetCustomerResponse, error) {
	customer, err := c.customerRepository.GetCustomerByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &dto.GetCustomerResponse{
		Name:   customer.Name,
		Points: customer.Points,
	}, nil
}

func (c *CustomerUseCaseImpl) CreateCustomer(ctx context.Context, customer *dto.CreateCustomerRequest) (*dto.CreateCustomerResponse, error) {
	convertedCustomer := &models.Customer{
		Name:   customer.Name,
		Points: customer.Points,
	}

	createdCustomer, err := c.customerRepository.CreateCustomer(ctx, convertedCustomer)
	if err != nil {
		return nil, err
	}

	return &dto.CreateCustomerResponse{
		Name:   createdCustomer.Name,
		Points: createdCustomer.Points,
	}, nil
}
