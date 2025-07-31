package usecase

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/internal/product/dto"
	httphelper "bsnack/app/pkg/http"
	"net/http"
	"strings"
	"time"

	"context"

	"gorm.io/gorm"
)

type ProductUseCaseImpl struct {
	productRepository interfaces.ProductRepository
}

func NewProductUseCase(productRepository interfaces.ProductRepository) interfaces.ProductUseCase {
	return &ProductUseCaseImpl{
		productRepository: productRepository,
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

func (p *ProductUseCaseImpl) GetProductsByManufactureDate(ctx context.Context, manufactureDate time.Time) (*[]models.Product, int64, error) {
	return p.productRepository.GetProductsByManufactureDate(ctx, manufactureDate)
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
