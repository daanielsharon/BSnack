package handler

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/product/dto"
	httphelper "bsnack/app/pkg/http"
	"bsnack/app/pkg/validation"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ProductHandlerImpl struct {
	ProductUseCase interfaces.ProductUseCase
}

func NewProductHandler(clientUseCase interfaces.ProductUseCase) interfaces.ProductHandler {
	return &ProductHandlerImpl{
		ProductUseCase: clientUseCase,
	}
}

func (p *ProductHandlerImpl) GetProductsByManufactureDate(w http.ResponseWriter, r *http.Request) {
	manufactureDate := r.URL.Query().Get("manufacture_date")
	manufactureDateParsed, err := time.Parse("2006-01-02", manufactureDate)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid manufacture date", nil)
		return
	}

	products, err := p.ProductUseCase.GetProductsByManufactureDate(r.Context(), manufactureDateParsed)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusInternalServerError, "Failed to get clients", nil)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Clients retrieved successfully", products)
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

	cases := cases.Title(language.Indonesian)

	product.Name = cases.String(product.Name)
	product.Type = cases.String(product.Type)
	product.Flavor = cases.String(product.Flavor)
	product.Size = cases.String(product.Size)

	createdProduct, err := p.ProductUseCase.CreateProduct(r.Context(), &product)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusInternalServerError, "Failed to create product", nil)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Product created successfully", createdProduct)
}
