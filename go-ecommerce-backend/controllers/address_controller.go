package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/services"
)

type AddressController struct {
	AddressService services.AddressService
}

func NewAddressController(service services.AddressService) *AddressController {
	return &AddressController{
		AddressService: service,
	}
}

func (ac *AddressController) GetAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}
	address, err := ac.AddressService.GetAddressByID(id)
	if err != nil {
		http.Error(w, "Address not found", http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(address)
}

func (ac *AddressController) CreateAddress(w http.ResponseWriter, r *http.Request, req *services.AddressRequest) {
	address, err := ac.AddressService.CreateAddress(req)

	if err != nil {
		http.Error(w, "Failed to create address", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address)
}

func (ac *AddressController) UpdateAddress(w http.ResponseWriter, r *http.Request, req *services.AddressRequest) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	address, err := ac.AddressService.UpdateAddress(id, req)

	if err != nil {
		http.Error(w, "Failed to update address", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(address)
}

func (ac *AddressController) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	err = ac.AddressService.DeleteAddress(id)

	if err != nil {
		http.Error(w, "Failed to delete address", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
