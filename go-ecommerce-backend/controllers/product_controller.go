package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/services"
)

type ProductController struct {
	ProductService services.ProductService
}

func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{
		ProductService: service,
	}
}

func (pc *ProductController) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := pc.ProductService.GetProductByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func (pc *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	orders, err := pc.ProductService.GetAllProducts()
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func (pc *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request, req *services.ProductRequest) {
	product, err := pc.ProductService.CreateProduct(req)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (pc *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request, req *services.ProductRequest) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := pc.ProductService.UpdateProduct(id, req)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func (pc *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = pc.ProductService.DeleteProduct(id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
