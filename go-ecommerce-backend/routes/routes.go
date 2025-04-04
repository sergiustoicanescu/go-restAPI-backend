package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/config"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/controllers"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/middlewares"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
)

func SetupRoutes(db *sql.DB, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	ctrls := controllers.NewControllers(db, cfg)

	setupPublicRoutes(router, ctrls.AuthController, ctrls.ProductController)

	v1 := router.PathPrefix("/v1").Subrouter()
	v1.Use(middlewares.JWTAuthMiddleware(cfg.JWTSecret))

	setupUserRoutes(v1, ctrls.UserController, ctrls.CustomerController)
	setupCustomerRoutes(v1, ctrls.CustomerController, ctrls.AddressController, ctrls.OrderController)
	setupAddressRoutes(v1, ctrls.AddressController)
	setupProductRoutes(v1, ctrls.ProductController)
	setupOrderRoutes(v1, ctrls.OrderController)

	return router
}

func setupOrderRoutes(v1 *mux.Router, orderController *controllers.OrderController) {
	ordersRoutes := v1.PathPrefix("/orders").Subrouter()
	ordersRoutes.Use(middlewares.OwnerOnlyMiddleware("id", orderController.OrderService.GetOwnerID))

	ordersRoutes.HandleFunc("/{id:[0-9]+}", orderController.GetOrder).Methods("GET")
	v1.HandleFunc("/orders", middlewares.ValidateBody(orderController.CreateOrder)).Methods("POST")
	ordersRoutes.HandleFunc("/{id:[0-9]+}", middlewares.ValidateBody(orderController.UpdateOrder)).Methods("PUT")
}

func setupProductRoutes(v1 *mux.Router, productController *controllers.ProductController) {
	productRoutes := v1.PathPrefix("/products").Subrouter()
	productRoutes.Use(middlewares.RoleAuthorizationMiddleware(string(models.RoleAdmin)))

	v1.HandleFunc("/products/{id:[0-9]+}", productController.GetProduct).Methods("GET")
	productRoutes.HandleFunc("", middlewares.ValidateBody(productController.CreateProduct)).Methods("POST")
	productRoutes.HandleFunc("/{id:[0-9]+}", middlewares.ValidateBody(productController.UpdateProduct)).Methods("PUT")
	productRoutes.HandleFunc("/{id:[0-9]+}", productController.DeleteProduct).Methods("DELETE")
}

func setupAddressRoutes(v1 *mux.Router, addressController *controllers.AddressController) {
	addressRoutes := v1.PathPrefix("/addresses").Subrouter()
	addressRoutes.Use(middlewares.OwnerOnlyMiddleware("id", addressController.AddressService.GetOwnerID))

	addressRoutes.HandleFunc("/{id:[0-9]+}", addressController.GetAddress).Methods("GET")
	v1.HandleFunc("/addresses", middlewares.ValidateBody(addressController.CreateAddress)).Methods("POST")
	addressRoutes.HandleFunc("/{id:[0-9]+}", middlewares.ValidateBody(addressController.UpdateAddress)).Methods("PUT")
	addressRoutes.HandleFunc("/{id:[0-9]+}", addressController.DeleteAddress).Methods("DELETE")
}

func setupCustomerRoutes(v1 *mux.Router, customerController *controllers.CustomerController, addressController *controllers.AddressController, orderController *controllers.OrderController) {
	customerRoutes := v1.PathPrefix("/customers").Subrouter()
	customerRoutes.Use(middlewares.OwnerOnlyMiddleware("id", customerController.CustomerService.GetOwnerID))

	customerRoutes.HandleFunc("/{id:[0-9]+}", customerController.GetCustomer).Methods("GET")
	v1.HandleFunc("/customers", middlewares.ValidateBody(customerController.CreateCustomer)).Methods("POST")
	customerRoutes.HandleFunc("/{id:[0-9]+}", middlewares.ValidateBody(customerController.UpdateCustomer)).Methods("PUT")
	customerRoutes.HandleFunc("/{id:[0-9]+}", customerController.DeleteCustomer).Methods("DELETE")

	customerRoutes.HandleFunc("/{id:[0-9]+}/addresses", addressController.GetAddressesByCustomerID).Methods("GET")
	customerRoutes.HandleFunc("/{id:[0-9]+}/orders", orderController.GetOrdersByCustomerID).Methods("GET")
}

func setupUserRoutes(v1 *mux.Router, userController *controllers.UserController, customerController *controllers.CustomerController) {
	userRoutes := v1.PathPrefix("/users").Subrouter()
	userRoutes.Use(middlewares.OwnerOnlyMiddleware("id", userController.UserService.GetOwnerID))

	userRoutes.HandleFunc("/{id:[0-9]+}", userController.GetUser).Methods("GET")
	userRoutes.HandleFunc("/{id:[0-9]+}", middlewares.ValidateBody(userController.UpdateUser)).Methods("PUT")
	userRoutes.HandleFunc("/{id:[0-9]+}/password", middlewares.ValidateBody(userController.UpdateUserPassword)).Methods("PATCH")
	userRoutes.HandleFunc("/{id:[0-9]+}", userController.DeleteUser).Methods("DELETE")

	userRoutes.HandleFunc("/{id:[0-9]+}/customer", customerController.GetCustomerByUserID).Methods("GET")
}

func setupPublicRoutes(router *mux.Router, authController *controllers.AuthController, productsController *controllers.ProductController) {
	router.HandleFunc("/v1/login", middlewares.ValidateBody(authController.Login)).Methods("POST")
	router.HandleFunc("/v1/register", middlewares.ValidateBody(authController.Register)).Methods("POST")
	router.HandleFunc("/v1/products", productsController.GetAllProducts).Methods("GET")
}
