package handler

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/product/dto"
	httphelper "bsnack/app/pkg/http"
	"bsnack/app/pkg/validation"
	"encoding/json"
	"net/http"
	"strings"
	"time"
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

	products, err := p.ProductUseCase.GetProductsByManufactureDate(r.Context(), manufactureDateParsed)
	if err != nil {
		httphelper.HandleError(w, err)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Products retrieved successfully", products)
}

func (p *ProductHandlerImpl) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid product data", nil)
		return
	}

	if err := validation.Validate.Struct(product); err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid product data", nil)
		return
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
