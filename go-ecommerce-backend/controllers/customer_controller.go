package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/services"
)

type CustomerController struct {
	CustomerService services.CustomerService
}

func NewCustomerController(service services.CustomerService) *CustomerController {
	return &CustomerController{
		CustomerService: service,
	}
}

func (cc *CustomerController) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customer, err := cc.CustomerService.GetCustomerByID(id)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func (cc *CustomerController) GetCustomerByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	customer, err := cc.CustomerService.GetCustomerByUserID(id)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func (cc *CustomerController) CreateCustomer(w http.ResponseWriter, r *http.Request, req *services.CustomerRequest) {
	customer, err := cc.CustomerService.CreateCustomer(req)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func (cc *CustomerController) UpdateCustomer(w http.ResponseWriter, r *http.Request, req *services.CustomerRequest) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customer, err := cc.CustomerService.UpdateCustomer(id, req)
	if err != nil {
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func (cc *CustomerController) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	err = cc.CustomerService.DeleteCustomer(id)
	if err != nil {
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
