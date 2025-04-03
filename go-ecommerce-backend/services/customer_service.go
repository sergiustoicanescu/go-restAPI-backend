package services

import (
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/repositories"
)

type CustomerService interface {
	GetCustomerByID(id int) (*models.Customer, error)
	CreateCustomer(req *CustomerRequest) (*models.Customer, error)
	UpdateCustomer(id int, req *CustomerRequest) (*models.Customer, error)
	DeleteCustomer(id int) error
	GetOwnerID(id int) (int, error)
}

type customerService struct {
	CustomerRepo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) CustomerService {
	return &customerService{
		CustomerRepo: repo,
	}
}

type CustomerRequest struct {
	UserID      int    `json:"user_id" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,max=15"`
}

func (cs *customerService) GetCustomerByID(id int) (*models.Customer, error) {
	return cs.CustomerRepo.GetByID(id)
}

func (cs *customerService) CreateCustomer(req *CustomerRequest) (*models.Customer, error) {
	customer := &models.Customer{
		UserID:      req.UserID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}
	err := cs.CustomerRepo.Create(customer)
	return customer, err
}

func (cs *customerService) UpdateCustomer(id int, req *CustomerRequest) (*models.Customer, error) {
	customer, err := cs.CustomerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	customer.FirstName = req.FirstName
	customer.LastName = req.LastName
	customer.PhoneNumber = req.PhoneNumber
	err = cs.CustomerRepo.Update(customer)
	return customer, err
}

func (cs *customerService) DeleteCustomer(id int) error {
	return cs.CustomerRepo.Delete(id)
}

func (cs *customerService) GetOwnerID(id int) (int, error) {
	customer, err := cs.GetCustomerByID(id)
	if err != nil {
		return 0, nil
	}
	return customer.UserID, err
}
