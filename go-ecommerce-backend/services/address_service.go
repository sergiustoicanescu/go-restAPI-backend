package services

import (
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/repositories"
)

type AddressService interface {
	GetAddressByID(id int) (*models.Address, error)
	CreateAddress(req *AddressRequest) (*models.Address, error)
	UpdateAddress(id int, req *AddressRequest) (*models.Address, error)
	DeleteAddress(id int) error
	GetOwnerID(id int) (int, error)
}

type addressService struct {
	AddressRepo repositories.AddressRepository
}

func NewAddressService(addressRepo repositories.AddressRepository) AddressService {
	return &addressService{
		AddressRepo: addressRepo,
	}
}

type AddressRequest struct {
	CustomerID    int    `json:"customer_id" validate:"required"`
	StreetAddress string `json:"street_address" validate:"required"`
	City          string `json:"city" validate:"required"`
	Country       string `json:"country" validate:"required"`
}

func (as *addressService) GetAddressByID(id int) (*models.Address, error) {
	return as.AddressRepo.GetByID(id)
}

func (as *addressService) CreateAddress(req *AddressRequest) (*models.Address, error) {
	address := &models.Address{
		CustomerID:    req.CustomerID,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		Country:       req.Country,
	}
	err := as.AddressRepo.Create(address)
	return address, err
}

func (as *addressService) UpdateAddress(id int, req *AddressRequest) (*models.Address, error) {
	address, err := as.AddressRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	address.StreetAddress = req.StreetAddress
	address.City = req.City
	address.Country = req.Country
	err = as.AddressRepo.Update(address)
	return address, err
}

func (as *addressService) DeleteAddress(id int) error {
	return as.AddressRepo.Delete(id)
}

func (as *addressService) GetOwnerID(id int) (int, error) {
	return as.AddressRepo.GetOwnerID(id)
}
