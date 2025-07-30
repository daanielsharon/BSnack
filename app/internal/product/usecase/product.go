package usecase

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/internal/product/dto"
	"time"

	"context"
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
	return p.productRepository.GetProductByName(ctx, name)
}

func (p *ProductUseCaseImpl) GetProductsByManufactureDate(ctx context.Context, manufactureDate time.Time) (*[]dto.GetProductResponse, error) {
	products, err := p.productRepository.GetProductsByManufactureDate(ctx, manufactureDate)
	if err != nil {
		return nil, err
	}

	var productResponse []dto.GetProductResponse
	for _, product := range *products {
		productResponse = append(productResponse, dto.GetProductResponse{
			Name:   product.Name,
			Type:   product.Type,
			Flavor: product.Flavor,
			Size:   product.Size,
			Price:  product.Price,
		})
	}
	return &productResponse, nil
}

func (p *ProductUseCaseImpl) CreateProduct(ctx context.Context, product *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
	convertedProduct := &models.Product{
		Name:            product.Name,
		Type:            product.Type,
		Flavor:          product.Flavor,
		Size:            product.Size,
		Price:           product.Price,
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
