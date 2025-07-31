package usecase

import (
	"bsnack/app/internal/customer/dto"
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/internal/shared"
	"bsnack/app/internal/validation"
	"bsnack/app/pkg/cache"
	httphelper "bsnack/app/pkg/http"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CustomerUseCaseImpl struct {
	customerRepository interfaces.CustomerRepository
	productUseCase     interfaces.ProductUseCase
	redisClient        *redis.Client
}

func NewCustomerUseCase(customerRepository interfaces.CustomerRepository, productUseCase interfaces.ProductUseCase, redisClient *redis.Client) interfaces.CustomerUseCase {
	return &CustomerUseCaseImpl{
		customerRepository: customerRepository,
		productUseCase:     productUseCase,
		redisClient:        redisClient,
	}
}

func (c *CustomerUseCaseImpl) GetCustomers(ctx context.Context) (*[]models.Customer, int64, error) {
	pg := shared.GetPagination(ctx)
	countKey := "customers:count"
	listKey := fmt.Sprintf("customers:page=%d:limit=%d", pg.Page, pg.PerPage)
	ttl := 5 * time.Minute

	listVal, err := c.redisClient.Get(ctx, listKey).Result()
	if err == nil {
		var customers []models.Customer
		err := json.Unmarshal([]byte(listVal), &customers)
		if err == nil {
			countVal, err := c.redisClient.Get(ctx, countKey).Result()
			if err == nil {
				total, _ := strconv.ParseInt(countVal, 10, 64)
				return &customers, total, nil
			}
		}
	}

	customers, total, err := c.customerRepository.GetCustomers(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = cache.SetJSON(ctx, c.redisClient, listKey, *customers, ttl)
	if err != nil {
		log.Printf("[WARN] Failed to set cache for list in get customers handler: %v", err)
	}

	err = cache.SetCount(ctx, c.redisClient, countKey, int64(len(*customers)), ttl)
	if err != nil {
		log.Printf("[WARN] Failed to set cache for count in get customers handler: %v", err)
	}

	return customers, total, nil
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
	if _, err := c.GetCustomerByName(ctx, customer.Name); err == nil {
		return nil, httphelper.NewAppError(http.StatusConflict, "Customer already exists")
	}
	convertedCustomer := &models.Customer{
		Name:   customer.Name,
		Points: customer.Points,
	}

	createdCustomer, err := c.customerRepository.CreateCustomer(ctx, convertedCustomer)
	if err != nil {
		return nil, err
	}

	customerPattern := "customers:*"
	err = cache.DeleteRedisKeysByPattern(ctx, c.redisClient, customerPattern)
	if err != nil {
		log.Printf("[WARN] Failed to delete cache for pattern %s in create customer handler: %v", customerPattern, err)
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

	productPattern := fmt.Sprintf("products:*:date=%s", product.GetDBManufactureDateInCorrectFormat())
	err = cache.DeleteRedisKeysByPattern(ctx, c.redisClient, productPattern)
	if err != nil {
		log.Printf("[WARN] Failed to delete cache for pattern %s in create point redemption handler: %v", productPattern, err)
	}

	customerPattern := "customers:*"
	err = cache.DeleteRedisKeysByPattern(ctx, c.redisClient, customerPattern)
	if err != nil {
		log.Printf("[WARN] Failed to delete cache for pattern %s in create point redemption handler: %v", customerPattern, err)
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
