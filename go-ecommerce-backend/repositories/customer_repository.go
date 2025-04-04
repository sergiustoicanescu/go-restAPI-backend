package repositories

import (
	"database/sql"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
)

type CustomerRepository interface {
	Create(customer *models.Customer) error
	GetByID(id int) (*models.Customer, error)
	GetByUserID(id int) (*models.Customer, error)
	Update(customer *models.Customer) error
	Delete(id int) error
}

type customerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{DB: db}
}

func (cr *customerRepository) Create(customer *models.Customer) error {
	query := "INSERT INTO customers (user_id, first_name, last_name, phone_number) VALUES ($1, $2, $3, $4) RETURNING id"
	return cr.DB.QueryRow(query, customer.UserID, customer.FirstName, customer.LastName, customer.PhoneNumber).Scan(&customer.ID)
}

func (cr *customerRepository) GetByID(id int) (*models.Customer, error) {
	customer := &models.Customer{}
	query := "SELECT id, user_id, first_name, last_name, phone_number FROM customers WHERE id = $1"
	err := cr.DB.QueryRow(query, id).Scan(&customer.ID, &customer.UserID, &customer.FirstName, &customer.LastName, &customer.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (cr *customerRepository) GetByUserID(id int) (*models.Customer, error) {
	customer := &models.Customer{}
	query := "SELECT id, user_id, first_name, last_name, phone_number FROM customers WHERE user_id = $1"
	err := cr.DB.QueryRow(query, id).Scan(&customer.ID, &customer.UserID, &customer.FirstName, &customer.LastName, &customer.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (cr *customerRepository) Update(customer *models.Customer) error {
	query := "UPDATE customers SET first_name = $1, last_name = $2, phone_number = $3 WHERE id = $4"
	_, err := cr.DB.Exec(query, customer.FirstName, customer.LastName, customer.PhoneNumber, customer.ID)
	return err
}

func (cr *customerRepository) Delete(id int) error {
	query := "DELETE FROM customers WHERE id = $1"
	_, err := cr.DB.Exec(query, id)
	return err
}
