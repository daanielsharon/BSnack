package repository

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/pkg/middleware"
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

func (p *ProductRepositoryImpl) GetProductByName(ctx context.Context, name string) (*models.Product, error) {
	var product *models.Product
	err := p.DB.WithContext(ctx).Where("name = ?", name).First(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductRepositoryImpl) GetProductsByManufactureDate(ctx context.Context, manufactureDate time.Time) (*[]models.Product, error) {
	pg := middleware.GetPagination(ctx)

	var products []models.Product
	var total int64

	db := p.DB.WithContext(ctx).Model(&models.Product{})
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	err := db.Offset(pg.Offset).Limit(pg.PerPage).Where("manufacture_date = ?", manufactureDate).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return &products, nil
}

func (p *ProductRepositoryImpl) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	return product, p.DB.WithContext(ctx).Model(&models.Product{}).Create(product).Error
}

func (p *ProductRepositoryImpl) DeductProductStock(ctx context.Context, productName string, quantity int) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		product, err := p.GetProductByName(ctx, productName)
		if err != nil {
			return err
		}

		product.Quantity -= quantity
		return tx.Model(&models.Product{}).Where("name = ?", productName).Update("quantity", product.Quantity).Error
	})
}
