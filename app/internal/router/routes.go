package router

import (
	v1 "bsnack/app/internal/api/v1"
	customer_handler "bsnack/app/internal/customer/handler"
	customer_repository "bsnack/app/internal/customer/repository"
	customer_usecase "bsnack/app/internal/customer/usecase"
	product_handler "bsnack/app/internal/product/handler"
	product_repository "bsnack/app/internal/product/repository"
	product_usecase "bsnack/app/internal/product/usecase"
	transaction_handler "bsnack/app/internal/transaction/handler"
	transaction_repository "bsnack/app/internal/transaction/repository"
	transaction_usecase "bsnack/app/internal/transaction/usecase"
	commonrouter "bsnack/app/pkg/router"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) chi.Router {
	r := commonrouter.New()

	productRepository := product_repository.NewProductRepository(DB)
	productUseCase := product_usecase.NewProductUseCase(productRepository, redisClient)
	productHandler := product_handler.NewProductHandler(productUseCase, redisClient)

	customerRepository := customer_repository.NewCustomerRepository(DB)
	customerUseCase := customer_usecase.NewCustomerUseCase(customerRepository, productUseCase, redisClient)
	customerHandler := customer_handler.NewCustomerHandler(customerUseCase)

	transactionRepository := transaction_repository.NewTransactionRepository(DB)
	transactionUseCase := transaction_usecase.NewTransactionUseCase(transactionRepository, customerUseCase, productUseCase, redisClient)
	transactionHandler := transaction_handler.NewTransactionHandler(transactionUseCase, productUseCase)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/v1", v1.Routes(v1.Router{
			CustomerHandler:    customerHandler,
			ProductHandler:     productHandler,
			TransactionHandler: transactionHandler,
		}))
	})

	return r
}
