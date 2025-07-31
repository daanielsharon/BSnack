package handler

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/product/dto"
	"bsnack/app/internal/shared"
	httphelper "bsnack/app/pkg/http"
	"bsnack/app/pkg/validation"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

type ProductHandlerImpl struct {
	ProductUseCase interfaces.ProductUseCase
}

func NewProductHandler(productUseCase interfaces.ProductUseCase) interfaces.ProductHandler {
	return &ProductHandlerImpl{
		ProductUseCase: productUseCase,
	}
}

func (p *ProductHandlerImpl) GetProductsByManufactureDate(w http.ResponseWriter, r *http.Request) {
	manufactureDate := r.URL.Query().Get("manufacture_date")
	if manufactureDate == "" {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Manufacture date is required", nil)
		return
	}

	manufactureDateParsed, err := time.Parse("2006-01-02", manufactureDate)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid manufacture date", nil)
		return
	}

	products, total, err := p.ProductUseCase.GetProductsByManufactureDate(r.Context(), manufactureDateParsed)
	if err != nil {
		httphelper.HandleError(w, err)
		return
	}

	productsResponse := make([]dto.GetProductResponse, len(*products))
	for i, product := range *products {
		productsResponse[i] = dto.GetProductResponse{
			Name:     product.Name,
			Type:     product.Type,
			Flavor:   product.Flavor,
			Size:     product.Size,
			Price:    product.Price,
			Quantity: product.Quantity,
		}
	}

	httphelper.JSONResponse(w, http.StatusOK, "Products retrieved successfully", shared.PaginatedResponse[dto.GetProductResponse]{
		Data:  productsResponse,
		Total: total,
	})
}

func (p *ProductHandlerImpl) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid product data", nil)
		return
	}

	if err := validation.Validate.Struct(product); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, fieldErr := range validationErrors {
				lowerCaseField := strings.ToLower(fieldErr.Field())
				switch fieldErr.Tag() {
				case "product_size":
					errors["size"] = "Product size must be one of: small, medium, or large"
				case "product_flavor":
					errors["flavor"] = "Invalid product flavor"
				case "product_type":
					errors["type"] = "Product type must be 'keripik pangsit'"
				case "YYYY-MM-DD_dateFormat":
					errors["manufacture_date"] = "Format date must be YYYY-MM-DD"
				case "product_price":
					errors["price"] = "Price does not match product size"
				default:
					errors[lowerCaseField] = fmt.Sprintf("Field validation failed on '%s' tag", fieldErr.Tag())
				}
			}
			httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid product data", errors)
			return
		}
	}

	product.Name = strings.ToLower(product.Name)
	product.Type = strings.ToLower(product.Type)
	product.Flavor = strings.ToLower(product.Flavor)
	product.Size = strings.ToLower(product.Size)

	createdProduct, err := p.ProductUseCase.CreateProduct(r.Context(), &product)
	if err != nil {
		httphelper.HandleError(w, err)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Product created successfully", createdProduct)
}
