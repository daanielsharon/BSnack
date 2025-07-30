package router

import (
	v1 "bsnack/app/internal/api/v1"
	customer_handler "bsnack/app/internal/customer/handler"
	customer_repository "bsnack/app/internal/customer/repository"
	customer_usecase "bsnack/app/internal/customer/usecase"
	product_handler "bsnack/app/internal/product/handler"
	product_repository "bsnack/app/internal/product/repository"
	product_usecase "bsnack/app/internal/product/usecase"
	commonrouter "bsnack/app/pkg/router"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) chi.Router {
	r := commonrouter.New()

	customerRepository := customer_repository.NewCustomerRepository(DB)
	customerUseCase := customer_usecase.NewCustomerUseCase(customerRepository)
	customerHandler := customer_handler.NewCustomerHandler(customerUseCase)

	productRepository := product_repository.NewProductRepository(DB)
	productUseCase := product_usecase.NewProductUseCase(productRepository)
	productHandler := product_handler.NewProductHandler(productUseCase)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/v1", v1.Routes(v1.Router{
			CustomerHandler: customerHandler,
			ProductHandler:  productHandler,
		}))
	})

	return r
}
