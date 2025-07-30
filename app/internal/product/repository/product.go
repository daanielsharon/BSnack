package repository

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"time"

	"context"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) interfaces.ProductRepository {
	return &ProductRepositoryImpl{DB: db}
}

func (p *ProductRepositoryImpl) GetProductsByManufactureDate(ctx context.Context, manufactureDate time.Time) (*[]models.Product, error) {
	var products []models.Product
	err := p.DB.WithContext(ctx).Find(&products, "manufacture_date = ?", manufactureDate).Error
	if err != nil {
		return nil, err
	}
	return &products, nil
}

func (p *ProductRepositoryImpl) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	return product, p.DB.WithContext(ctx).Model(&models.Product{}).Create(product).Error
}
