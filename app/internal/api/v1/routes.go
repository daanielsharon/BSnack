package v1

import (
	"bsnack/app/internal/interfaces"
	"bsnack/app/pkg/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	CustomerHandler    interfaces.CustomerHandler
	ProductHandler     interfaces.ProductHandler
	TransactionHandler interfaces.TransactionHandler
}

func Routes(router Router) http.Handler {
	r := chi.NewRouter()

	r.Route("/customers", func(r chi.Router) {
		r.With(middleware.PaginationMiddleware).Get("/", router.CustomerHandler.GetCustomers)
		r.Route("/{name}", func(r chi.Router) {
			r.Get("/", router.CustomerHandler.GetCustomerByName)
			r.Post("/point-redemption", router.CustomerHandler.CreatePointRedemption)
		})
		r.Post("/", router.CustomerHandler.CreateCustomer)
	})

	r.Route("/products", func(r chi.Router) {
		r.With(middleware.PaginationMiddleware).Get("/", router.ProductHandler.GetProductsByManufactureDate)
		r.Post("/", router.ProductHandler.CreateProduct)
	})

	r.Route("/transactions", func(r chi.Router) {
		r.With(middleware.PaginationMiddleware).Get("/", router.TransactionHandler.GetTransactions)
		r.Post("/", router.TransactionHandler.CreateTransaction)
		r.Get("/{id}", router.TransactionHandler.GetTransactionById)
	})

	return r
}
