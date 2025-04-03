package repositories

import (
	"database/sql"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (ur *userRepository) Create(user *models.User) error {
	query := "INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id"
	return ur.DB.QueryRow(query, user.Email, user.Password, user.Role).Scan(&user.ID)
}

func (ur *userRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, email, password, role, created_at FROM users WHERE id = $1"
	err := ur.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, email, password, role, created_at FROM users WHERE email = $1"
	err := ur.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) Update(user *models.User) error {
	query := "UPDATE users SET email = $1, password = $2, role = $3 WHERE id = $4"
	_, err := ur.DB.Exec(query, user.Email, user.Password, user.Role, user.ID)
	return err
}

func (ur *userRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := ur.DB.Exec(query, id)
	return err
}
