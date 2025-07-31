package handler

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/internal/product/dto"
	"bsnack/app/internal/shared"
	"bsnack/app/pkg/cache"
	httphelper "bsnack/app/pkg/http"
	"bsnack/app/pkg/validation"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
)

type ProductHandlerImpl struct {
	ProductUseCase interfaces.ProductUseCase
	RedisClient    *redis.Client
}

func NewProductHandler(productUseCase interfaces.ProductUseCase, redisClient *redis.Client) interfaces.ProductHandler {
	return &ProductHandlerImpl{
		ProductUseCase: productUseCase,
		RedisClient:    redisClient,
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

	pg := shared.GetPagination(r.Context())
	countKey := fmt.Sprintf("products:count:date=%s", manufactureDate)
	listKey := fmt.Sprintf("products:page=%d:limit=%d:date=%s", pg.Page, pg.PerPage, manufactureDate)
	ttl := 10 * time.Minute

	listVal, err := p.RedisClient.Get(r.Context(), listKey).Result()
	if err == nil {
		var productsResponse []dto.GetProductResponse
		err := json.Unmarshal([]byte(listVal), &productsResponse)
		if err == nil {
			countVal, err := p.RedisClient.Get(r.Context(), countKey).Result()
			if err == nil {
				total, _ := strconv.ParseInt(countVal, 10, 64)
				httphelper.JSONResponse(w, http.StatusOK, "Products retrieved successfully", shared.PaginatedResponse[dto.GetProductResponse]{
					Data:  productsResponse,
					Total: total,
				})
				return
			}
		}
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

	err = cache.SetJSON(r.Context(), p.RedisClient, listKey, productsResponse, ttl)
	if err != nil {
		log.Printf("[WARN] Failed to set cache for list in get products by manufacture date handler: %v", err)
		return
	}
	err = cache.SetCount(r.Context(), p.RedisClient, countKey, total, ttl)
	if err != nil {
		log.Printf("[WARN] Failed to set cache for count in get products by manufacture date handler: %v", err)
		return
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

	dateStr := product.ManufactureDate
	pattern := fmt.Sprintf("products:*:date=%s", dateStr)

	err = cache.DeleteRedisKeysByPattern(r.Context(), p.RedisClient, pattern)
	if err != nil {
		log.Printf("[WARN] Failed to delete cache for pattern %s in create product handler: %v", pattern, err)
	}

	httphelper.JSONResponse(w, http.StatusOK, "Product created successfully", createdProduct)
}
