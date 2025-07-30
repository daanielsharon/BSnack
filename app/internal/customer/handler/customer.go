package handler

import (
	"bsnack/app/internal/customer/dto"
	"bsnack/app/internal/interfaces"
	httphelper "bsnack/app/pkg/http"
	"bsnack/app/pkg/validation"
	"encoding/json"
	"net/http"
)

type CustomerHandlerImpl struct {
	CustomerUseCase interfaces.CustomerUseCase
}

func NewCustomerHandler(customerUseCase interfaces.CustomerUseCase) interfaces.CustomerHandler {
	return &CustomerHandlerImpl{
		CustomerUseCase: customerUseCase,
	}
}

func (c *CustomerHandlerImpl) GetCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := c.CustomerUseCase.GetCustomers(r.Context())
	if err != nil {
		httphelper.JSONResponse(w, http.StatusInternalServerError, "Failed to get customers", nil)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Customers retrieved successfully", customers)
}

func (c *CustomerHandlerImpl) GetCustomerByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	customer, err := c.CustomerUseCase.GetCustomerByName(r.Context(), name)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusInternalServerError, "Failed to get customer", nil)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Customer retrieved successfully", customer)
}

func (c *CustomerHandlerImpl) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer dto.CreateCustomerRequest
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid customer data", nil)
		return
	}

	if err := validation.Validate.Struct(customer); err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid customer data", nil)
		return
	}

	createdCustomer, err := c.CustomerUseCase.CreateCustomer(r.Context(), &customer)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusInternalServerError, "Failed to create customer", nil)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Customer created successfully", createdCustomer)
}
