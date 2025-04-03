package controllers

import (
	"database/sql"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/config"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/repositories"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/services"
)

type AllControllers struct {
	AuthController     *AuthController
	UserController     *UserController
	CustomerController *CustomerController
	AddressController  *AddressController
	ProductController  *ProductController
	OrderController    *OrderController
}

func NewControllers(db *sql.DB, cfg *config.Config) *AllControllers {
	userRepo := repositories.NewUserRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	addressRepo := repositories.NewAddressRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo)
	customerService := services.NewCustomerService(customerRepo)
	addressService := services.NewAddressService(addressRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo)

	return &AllControllers{
		AuthController:     NewAuthController(authService),
		UserController:     NewUserController(userService),
		CustomerController: NewCustomerController(customerService),
		AddressController:  NewAddressController(addressService),
		ProductController:  NewProductController(productService),
		OrderController:    NewOrderController(orderService),
	}
}
