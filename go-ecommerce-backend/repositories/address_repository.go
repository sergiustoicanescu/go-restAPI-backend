package repositories

import (
	"database/sql"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
)

type AddressRepository interface {
	Create(address *models.Address) error
	GetByID(id int) (*models.Address, error)
	GetByCustomerID(id int) ([]*models.Address, error)
	Update(address *models.Address) error
	Delete(id int) error
	GetOwnerID(id int) (int, error)
}

type addressRepository struct {
	DB *sql.DB
}

func NewAddressRepository(db *sql.DB) AddressRepository {
	return &addressRepository{DB: db}
}

func (ar *addressRepository) Create(address *models.Address) error {
	query := "INSERT INTO addresses (customer_id, street_address, city, country) VALUES ($1, $2, $3, $4) RETURNING id"
	return ar.DB.QueryRow(query, address.CustomerID, address.StreetAddress, address.City, address.Country).Scan(&address.ID)
}

func (ar *addressRepository) GetByID(id int) (*models.Address, error) {
	address := &models.Address{}
	query := "SELECT id, customer_id, street_address, city, country FROM addresses WHERE id = $1"
	err := ar.DB.QueryRow(query, id).Scan(&address.ID, &address.CustomerID, &address.StreetAddress, &address.City, &address.Country)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (ar *addressRepository) GetByCustomerID(id int) ([]*models.Address, error) {
	query := "SELECT id, customer_id, street_address, city, country FROM addresses WHERE customer_id = $1"
	rows, err := ar.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*models.Address
	for rows.Next() {
		var address = new(models.Address)
		if err := rows.Scan(&address.ID, &address.CustomerID, &address.StreetAddress, &address.City, &address.Country); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func (ar *addressRepository) Update(address *models.Address) error {
	query := "UPDATE addresses SET street_address = $1, city = $2, country = $3 WHERE id = $4"
	_, err := ar.DB.Exec(query, address.StreetAddress, address.City, address.Country, address.ID)
	return err
}

func (ar *addressRepository) Delete(id int) error {
	query := "DELETE FROM addresses WHERE id = $1"
	_, err := ar.DB.Exec(query, id)
	return err
}

func (ar *addressRepository) GetOwnerID(id int) (int, error) {
	var userID int
	query := "SELECT c.user_id FROM addresses a JOIN customers c ON c.id = a.customer_id WHERE a.id = $1"
	err := ar.DB.QueryRow(query, id).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, err
}
