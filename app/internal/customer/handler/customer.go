package handler

import (
	"bsnack/app/internal/customer/dto"
	"bsnack/app/internal/interfaces"
	httphelper "bsnack/app/pkg/http"
	"bsnack/app/pkg/validation"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CustomerHandlerImpl struct {
	customerUseCase interfaces.CustomerUseCase
}

func NewCustomerHandler(customerUseCase interfaces.CustomerUseCase) interfaces.CustomerHandler {
	return &CustomerHandlerImpl{
		customerUseCase: customerUseCase,
	}
}

func (c *CustomerHandlerImpl) GetCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := c.customerUseCase.GetCustomers(r.Context())
	if err != nil {
		httphelper.HandleError(w, err)
		return
	}

	customersResponse := make([]dto.GetCustomerResponse, len(*customers))
	for i, customer := range *customers {
		customersResponse[i] = dto.GetCustomerResponse{
			Name:   customer.Name,
			Points: customer.Points,
		}
	}

	httphelper.JSONResponse(w, http.StatusOK, "Customers retrieved successfully", customersResponse)
}

func (c *CustomerHandlerImpl) GetCustomerByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	customer, err := c.customerUseCase.GetCustomerByName(r.Context(), name)
	if err != nil {
		httphelper.HandleError(w, err)
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

	createdCustomer, err := c.customerUseCase.CreateCustomer(r.Context(), &customer)
	if err != nil {
		httphelper.HandleError(w, err)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Customer created successfully", createdCustomer)
}

func (c *CustomerHandlerImpl) CreatePointRedemption(w http.ResponseWriter, r *http.Request) {
	customerName := chi.URLParam(r, "name")
	var pointRedemption dto.CreatePointRedemptionRequest
	err := json.NewDecoder(r.Body).Decode(&pointRedemption)
	if err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid point redemption data", nil)
		return
	}

	fmt.Println("customerName", customerName)
	fmt.Println("pointRedemption", pointRedemption)

	if err := validation.Validate.Struct(pointRedemption); err != nil {
		httphelper.JSONResponse(w, http.StatusBadRequest, "Invalid point redemption data", nil)
		return
	}

	createdPointRedemption, err := c.customerUseCase.CreatePointRedemption(r.Context(), customerName, &pointRedemption)
	if err != nil {
		httphelper.HandleError(w, err)
		return
	}

	httphelper.JSONResponse(w, http.StatusOK, "Point redemption created successfully", createdPointRedemption)
}
