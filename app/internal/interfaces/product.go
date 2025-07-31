package interfaces

import (
	"bsnack/app/internal/models"
	"bsnack/app/internal/product/dto"
	"context"
	"net/http"
	"time"
)

type ProductHandler interface {
	GetProductsByManufactureDate(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
}

type ProductUseCase interface {
	GetProductsByManufactureDate(ctx context.Context, manufactureDate time.Time) (*[]models.Product, error)
	CreateProduct(ctx context.Context, product *dto.CreateProductRequest) (*dto.CreateProductResponse, error)
	GetProductByName(ctx context.Context, name string) (*models.Product, error)
}

type ProductRepository interface {
	GetProductsByManufactureDate(ctx context.Context, manufactureDate time.Time) (*[]models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProductByName(ctx context.Context, name string) (*models.Product, error)
}
