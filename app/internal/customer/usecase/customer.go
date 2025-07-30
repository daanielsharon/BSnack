package usecase

import (
	"bsnack/app/internal/customer/dto"
	"bsnack/app/internal/customer/services"
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"strings"

	"context"
)

type CustomerUseCaseImpl struct {
	customerRepository interfaces.CustomerRepository
	productUseCase     interfaces.ProductUseCase
}

func NewCustomerUseCase(customerRepository interfaces.CustomerRepository, productUseCase interfaces.ProductUseCase) interfaces.CustomerUseCase {
	return &CustomerUseCaseImpl{
		customerRepository: customerRepository,
		productUseCase:     productUseCase,
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

func (c *CustomerUseCaseImpl) CreatePointRedemption(ctx context.Context, pointRedemption *dto.CreatePointRedemptionRequest) (*dto.CreatePointRedemptionResponse, error) {
	product, err := c.productUseCase.GetProductByName(ctx, pointRedemption.ProductName)
	if err != nil {
		return nil, err
	}

	err = services.ValidateProductExists(product)
	if err != nil {
		return nil, err
	}

	err = services.ValidateProductStockEnough(product, pointRedemption.Quantity)
	if err != nil {
		return nil, err
	}

	customer, err := c.customerRepository.GetCustomerByName(ctx, pointRedemption.CustomerName)
	if err != nil {
		return nil, err
	}

	productSize := strings.ToLower(pointRedemption.ProductSize)
	var sizePointMap = map[string]int{
		"small":  200,
		"medium": 300,
		"large":  500,
	}

	pointRequired := sizePointMap[productSize] * pointRedemption.Quantity
	err = services.ValidateEnoughPoint(customer.Points, pointRequired)
	if err != nil {
		return nil, err
	}

	convertedPointRedemption := &models.PointRedemption{
		CustomerName:  customer.Name,
		ProductName:   product.Name,
		ProductSize:   productSize,
		ProductType:   product.Type,
		ProductFlavor: product.Flavor,
		Quantity:      pointRedemption.Quantity,
		PointRequired: pointRequired,
	}

	createdPointRedemption, err := c.customerRepository.CreatePointRedemption(ctx, convertedPointRedemption)
	if err != nil {
		return nil, err
	}

	updatedCustomer, err := c.customerRepository.UpdateCustomerPoints(ctx, customer)
	if err != nil {
		return nil, err
	}

	return &dto.CreatePointRedemptionResponse{
		ProductName:   product.Name,
		ProductSize:   productSize,
		ProductType:   product.Type,
		ProductFlavor: product.Flavor,
		PointBefore:   customer.Points,
		PointRequired: pointRequired,
		PointAfter:    updatedCustomer.Points,
		RedeemedAt:    createdPointRedemption.RedeemedAt,
	}, nil
}
