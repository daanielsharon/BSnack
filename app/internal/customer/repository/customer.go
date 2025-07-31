package repository

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"
	"bsnack/app/internal/shared"

	"context"

	"gorm.io/gorm"
)

type CustomerRepositoryImpl struct {
	DB *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) interfaces.CustomerRepository {
	return &CustomerRepositoryImpl{DB: db}
}

func (c *CustomerRepositoryImpl) GetCustomers(ctx context.Context) (*[]models.Customer, int64, error) {
	pg := shared.GetPagination(ctx)
	var customers []models.Customer
	var total int64

	db := c.DB.WithContext(ctx).Model(&models.Customer{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(pg.Offset).Limit(pg.PerPage).Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}
	return &customers, total, nil
}

func (c *CustomerRepositoryImpl) GetCustomerByName(ctx context.Context, name string) (*models.Customer, error) {
	var customer *models.Customer
	err := c.DB.WithContext(ctx).Where("name = ?", name).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c *CustomerRepositoryImpl) CreateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	return customer, c.DB.WithContext(ctx).Model(&models.Customer{}).Create(customer).Error
}

func (c *CustomerRepositoryImpl) AddCustomerPoints(ctx context.Context, customerName string, points int) error {
	return c.DB.Transaction(func(tx *gorm.DB) error {
		customer, err := c.GetCustomerByName(ctx, customerName)
		if err != nil {
			return err
		}

		customer.Points += points
		return tx.Model(&models.Customer{}).Where("name = ?", customerName).Update("points", customer.Points).Error
	})
}

func (c *CustomerRepositoryImpl) DeductCustomerPoints(ctx context.Context, customerName string, points int) error {
	return c.DB.Transaction(func(tx *gorm.DB) error {
		customer, err := c.GetCustomerByName(ctx, customerName)
		if err != nil {
			return err
		}

		customer.Points -= points
		return tx.Model(&models.Customer{}).Where("name = ?", customerName).Update("points", customer.Points).Error
	})
}

func (c *CustomerRepositoryImpl) CreatePointRedemption(ctx context.Context, pointRedemption *models.PointRedemption) (*models.PointRedemption, error) {
	return pointRedemption, c.DB.WithContext(ctx).Model(&models.PointRedemption{}).Create(pointRedemption).Error
}
