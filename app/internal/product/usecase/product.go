package usecase

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/internal/product/dto"
	"bsnack/app/internal/shared"
	"bsnack/app/pkg/cache"
	httphelper "bsnack/app/pkg/http"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ProductUseCaseImpl struct {
	productRepository interfaces.ProductRepository
	redisClient       *redis.Client
}

func NewProductUseCase(productRepository interfaces.ProductRepository, redisClient *redis.Client) interfaces.ProductUseCase {
	return &ProductUseCaseImpl{
		productRepository: productRepository,
		redisClient:       redisClient,
	}
}

func (p *ProductUseCaseImpl) GetProductByName(ctx context.Context, name string) (*models.Product, error) {
	product, err := p.productRepository.GetProductByName(ctx, strings.ToLower(name))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, httphelper.NewAppError(http.StatusNotFound, "Product not found")
		}
		return nil, err
	}
	return product, nil
}

func (p *ProductUseCaseImpl) GetProductsByManufactureDate(ctx context.Context, manufactureDate string) (*[]models.Product, int64, error) {
	pg := shared.GetPagination(ctx)
	countKey := fmt.Sprintf("products:count:date=%s", manufactureDate)
	listKey := fmt.Sprintf("products:page=%d:limit=%d:date=%s", pg.Page, pg.PerPage, manufactureDate)
	ttl := 10 * time.Minute

	listVal, err := p.redisClient.Get(ctx, listKey).Result()
	if err == nil {
		var productsResponse []models.Product
		err := json.Unmarshal([]byte(listVal), &productsResponse)
		if err == nil {
			countVal, err := p.redisClient.Get(ctx, countKey).Result()
			if err == nil {
				total, _ := strconv.ParseInt(countVal, 10, 64)
				return &productsResponse, total, nil
			}
		}
	}

	products, total, err := p.productRepository.GetProductsByManufactureDate(ctx, manufactureDate)
	if err != nil {
		return nil, 0, err
	}

	err = cache.SetJSON(ctx, p.redisClient, listKey, *products, ttl)
	if err != nil {
		log.Printf("[WARN] Failed to set cache for list in get products by manufacture date handler: %v", err)
	}
	err = cache.SetCount(ctx, p.redisClient, countKey, int64(len(*products)), ttl)
	if err != nil {
		log.Printf("[WARN] Failed to set cache for count in get products by manufacture date handler: %v", err)
	}

	return products, total, nil
}

func (p *ProductUseCaseImpl) CreateProduct(ctx context.Context, product *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
	convertedProduct := &models.Product{
		Name:            product.Name,
		Type:            product.Type,
		Flavor:          product.Flavor,
		Size:            product.Size,
		Price:           product.Price,
		Quantity:        product.Quantity,
		ManufactureDate: product.ManufactureDate,
	}

	createdProduct, err := p.productRepository.CreateProduct(ctx, convertedProduct)
	if err != nil {
		return nil, err
	}

	pattern := fmt.Sprintf("products:*:date=%s", product.ManufactureDate)
	err = cache.DeleteRedisKeysByPattern(ctx, p.redisClient, pattern)
	if err != nil {
		log.Printf("[WARN] Failed to delete cache for pattern %s in create product handler: %v", pattern, err)
	}

	return &dto.CreateProductResponse{
		Name:            createdProduct.Name,
		Type:            createdProduct.Type,
		Flavor:          createdProduct.Flavor,
		Size:            createdProduct.Size,
		Price:           createdProduct.Price,
		ManufactureDate: createdProduct.ManufactureDate,
		CreatedAt:       createdProduct.CreatedAt,
	}, nil
}

func (p *ProductUseCaseImpl) DeductProductStock(ctx context.Context, productName string, quantity int) error {
	return p.productRepository.DeductProductStock(ctx, productName, quantity)
}
