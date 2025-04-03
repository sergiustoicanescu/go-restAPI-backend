# Go REST API Backend

## Table of Contents
- [Project Overview](#project-overview)
- [Project Structure](#project-structure)
- [Setup Instructions](#setup-instructions)
  - [Database Setup](#database-setup)
  - [Environment Variables](#environment-variables)
  - [Running the Application](#running-the-application)
- [Available Endpoints](#available-endpoints)

## Project Overview

This project is a small e-commerce backend that provides endpoints for managing users, customers, addresses, products, and orders. Key features include:

- **Authentication & Authorization:**  
  - JWT-based authentication.
  - Role-based access control (e.g., customer, admin).

- **REST API with CRUD Operations:**  
  Endpoints to create, read, update, and delete resources such as users, products, orders, and customers.

- **Input Validation:**  
  Utilizes [go-playground/validator](https://github.com/go-playground/validator) for input validation.

- **SQL Database Integration:**  
  * Direct SQL queries using Go’s `database/sql` package with migration handling via [golang-migrate](https://github.com/golang-migrate/migrate). 
  * No ORM is used.

- **Modular Architecture:**  
  Clear separation of concerns across across all application layers.

---

## Project Structure

```
go-ecommerce-backend/
├── config          Loads configuration from environment variables.
├── controllers     Handles HTTP requests and responses.
├── middlewares     Implements authentication, authorization, and request interceptors.
├── migrations      Contains SQL migration files for managing the database schema.
├── models          Defines domain models and data structures.
├── repositories    Implements data access using raw SQL queries.
├── routes          Sets up HTTP routes and middleware chaining.
├── services        Orchestrates repository interactions.
├── unit_tests      Contains unit tests for the auth and order controllers.
├── utils           Provides utility functions.
└── main.go         Application entry point.
```
---

## Setup Instructions

### Database Setup
Create two databases for development and testing:

* `ecommerce` for the production/development environment.
* `ecommerce_test` for testing.

### Environment Variables

Create a `.env` file in the project root with these variables:

```env
DB_PROD=postgres://username:yourpassword@localhost:port/ecommerce?sslmode=disable
DB_TEST=postgres://username:yourpassword@localhost:port/ecommerce_test?sslmode=disable
JWT_SECRET=your-secret-key
```

### Running the Application

```
go run main.go
```

The API will be accessible at the port specified by the PORT environment variable (default is 8080).

## Available Endpoints
### Public Endpoints:

* `POST /v1/login` – Authenticate user and return a JWT token.

* `POST /v1/register` – Register a new user.

* `GET /v1/products` – Retrieve a list of all products.

### Protected Endpoints (JWT Required):

#### User Endpoints:

* `GET /v1/users/{id}` – Retrieve user details.

* `PUT /v1/users/{id}` – Update user details.

* `PATCH /v1/users/{id}/password` – Update user password.

* `DELETE /v1/users/{id}` – Delete a user.

#### Customer Endpoints:

* `GET /v1/customers/{id}` – Retrieve customer details.

* `POST /v1/customers` – Create a new customer.

* `PUT /v1/customers/{id}` – Update customer details.

* `DELETE /v1/customers/{id}` – Delete a customer.

#### Address Endpoints:

* `GET /v1/addresses/{id}` – Retrieve address details.

* `POST /v1/addresses` – Create a new address.

* `PUT /v1/addresses/{id}` – Update an address.

* `DELETE /v1/addresses/{id}` – Delete an address.

#### Product Endpoints (Admin Only):

* `POST /v1/products` – Create a new product.

* `PUT /v1/products/{id}` – Update product details.

* `DELETE /v1/products/{id}` – Delete a product.

#### Order Endpoints:

* `GET /v1/orders/{id}` – Retrieve order details.

* `POST /v1/orders` – Create a new order.

* `PUT /v1/orders/{id}` – Update order details.