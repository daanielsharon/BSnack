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
	GetProductsByManufactureDate(ctx context.Context, manufactureDate time.Time) (*[]dto.GetProductResponse, error)
	CreateProduct(ctx context.Context, product *dto.CreateProductRequest) (*dto.CreateProductResponse, error)
}

type ProductRepository interface {
	GetProductsByManufactureDate(ctx context.Context, manufactureDate time.Time) (*[]models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
}
