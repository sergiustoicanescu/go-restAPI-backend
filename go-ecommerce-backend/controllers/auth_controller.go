package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/services"
)

type AuthController struct {
	AuthService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request, req *services.LoginRequest) {
	user, token, err := a.AuthService.Login(req)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	response := map[string]any{
		"user":  user,
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (a *AuthController) Register(w http.ResponseWriter, r *http.Request, req *services.RegisterRequest) {
	err := a.AuthService.Register(req)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
