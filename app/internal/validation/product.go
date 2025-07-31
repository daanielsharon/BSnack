package validation

import (
	"bsnack/app/internal/models"
	httphelper "bsnack/app/pkg/http"
	"net/http"
)

func ValidateProductExists(product *models.Product) error {
	if product == nil {
		return httphelper.NewAppError(http.StatusNotFound, "Product not found")
	}
	return nil
}

func ValidateProductStockEnough(product *models.Product, quantity int) error {
	if quantity > product.Quantity {
		return httphelper.NewAppError(http.StatusBadRequest, "Product stock not enough")
	}
	return nil
}

func ValidateSameProduct(product *models.Product, productInput *models.Product) error {
	if product.Name != productInput.Name || product.Flavor != productInput.Flavor || product.Size != productInput.Size {
		return httphelper.NewAppError(http.StatusBadRequest, "Product not found")
	}
	return nil
}
