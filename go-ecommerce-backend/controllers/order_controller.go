package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/services"
)

type OrderController struct {
	OrderService services.OrderService
}

func NewOrderController(service services.OrderService) *OrderController {
	return &OrderController{
		OrderService: service,
	}
}

func (oc *OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := oc.OrderService.GetOrderByID(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func (oc *OrderController) GetOrdersByCustomerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	order, err := oc.OrderService.GetOrdersByCustomerID(id)
	if err != nil {
		http.Error(w, "Orders for customer not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func (oc *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request, req *services.OrderRequest) {
	order, err := oc.OrderService.CreateOrder(req)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (oc *OrderController) UpdateOrder(w http.ResponseWriter, r *http.Request, req *services.OrderRequest) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := oc.OrderService.UpdateOrder(id, req)
	if err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(order)
}
