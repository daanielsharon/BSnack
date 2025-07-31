package usecase

import (
	"bsnack/app/internal/customer/dto"
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/internal/validation"
	httphelper "bsnack/app/pkg/http"
	"context"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
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

func (c *CustomerUseCaseImpl) GetCustomers(ctx context.Context) (*[]models.Customer, error) {
	return c.customerRepository.GetCustomers(ctx)
}

func (c *CustomerUseCaseImpl) GetCustomerByName(ctx context.Context, name string) (*dto.GetCustomerResponse, error) {
	customer, err := c.customerRepository.GetCustomerByName(ctx, name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, httphelper.NewAppError(http.StatusNotFound, "Customer not found")
		}
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

func (c *CustomerUseCaseImpl) DeductCustomerPoint(ctx context.Context, customerName string, pointRequired int) (*dto.GetCustomerResponse, error) {
	if err := c.customerRepository.DeductCustomerPoints(ctx, customerName, pointRequired); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, httphelper.NewAppError(http.StatusNotFound, "Customer not found")
		}

		return nil, err
	}

	return c.GetCustomerByName(ctx, customerName)
}

func (c *CustomerUseCaseImpl) AddCustomerPoint(ctx context.Context, customerName string, points int) (*dto.GetCustomerResponse, error) {
	if err := c.customerRepository.AddCustomerPoints(ctx, customerName, points); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, httphelper.NewAppError(http.StatusNotFound, "Customer not found")
		}

		return nil, err
	}

	return c.GetCustomerByName(ctx, customerName)
}

func (c *CustomerUseCaseImpl) CreatePointRedemption(ctx context.Context, customerName string, pointRedemption *dto.CreatePointRedemptionRequest) (*dto.CreatePointRedemptionResponse, error) {
	productName := strings.ToLower(pointRedemption.ProductName)
	productSize := strings.ToLower(pointRedemption.ProductSize)
	productFlavor := strings.ToLower(pointRedemption.ProductFlavor)

	product, err := c.productUseCase.GetProductByName(ctx, productName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, httphelper.NewAppError(http.StatusNotFound, "Product not found")
		}

		return nil, err
	}

	err = validation.ValidateProductExists(product)
	if err != nil {
		return nil, err
	}

	err = validation.ValidateProductStockEnough(product, pointRedemption.Quantity)
	if err != nil {
		return nil, err
	}

	customer, err := c.customerRepository.GetCustomerByName(ctx, customerName)
	if err != nil {
		return nil, err
	}

	err = validation.ValidateSameProduct(product, &models.Product{
		Name:   productName,
		Size:   productSize,
		Flavor: productFlavor,
	})
	if err != nil {
		return nil, err
	}

	var sizePointMap = map[string]int{
		"small":  200,
		"medium": 300,
		"large":  500,
	}

	pointRequired := sizePointMap[productSize] * pointRedemption.Quantity
	err = validation.ValidateEnoughPoint(customer.Points, pointRequired)
	if err != nil {
		return nil, err
	}

	convertedPointRedemption := &models.PointRedemption{
		CustomerName:  customer.Name,
		ProductName:   productName,
		ProductSize:   productSize,
		ProductFlavor: productFlavor,
		Quantity:      pointRedemption.Quantity,
		PointRequired: pointRequired,
		RedeemedAt:    time.Now(),
	}

	createdPointRedemption, err := c.customerRepository.CreatePointRedemption(ctx, convertedPointRedemption)
	if err != nil {
		return nil, err
	}

	updatedCustomer, err := c.DeductCustomerPoint(ctx, customer.Name, pointRequired)
	if err != nil {
		return nil, err
	}

	if err := c.productUseCase.DeductProductStock(ctx, productName, pointRedemption.Quantity); err != nil {
		return nil, err
	}

	return &dto.CreatePointRedemptionResponse{
		ProductName:   productName,
		ProductSize:   productSize,
		ProductFlavor: productFlavor,
		PointBefore:   customer.Points,
		PointRequired: pointRequired,
		PointAfter:    updatedCustomer.Points,
		RedeemedAt:    createdPointRedemption.RedeemedAt,
	}, nil
}
