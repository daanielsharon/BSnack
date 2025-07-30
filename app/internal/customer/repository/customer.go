package repository

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/models"

	"context"

	"gorm.io/gorm"
)

type CustomerRepositoryImpl struct {
	DB *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) interfaces.CustomerRepository {
	return &CustomerRepositoryImpl{DB: db}
}

func (c *CustomerRepositoryImpl) GetCustomers(ctx context.Context) (*[]models.Customer, error) {
	var customers []models.Customer
	err := c.DB.WithContext(ctx).Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return &customers, nil
}

func (c *CustomerRepositoryImpl) GetCustomerByName(ctx context.Context, name string) (*models.Customer, error) {
	var customer models.Customer
	err := c.DB.WithContext(ctx).Find(&customer, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (c *CustomerRepositoryImpl) CreateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	return customer, c.DB.WithContext(ctx).Model(&models.Customer{}).Create(customer).Error
}
