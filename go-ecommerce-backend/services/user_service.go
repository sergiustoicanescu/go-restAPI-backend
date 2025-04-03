package services

import (
	"errors"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUserByID(id int) (*models.User, error)
	UpdateUser(id int, userReq *UserRequest) (*models.User, error)
	UpdateUserPassword(id int, req *UserRequest) (*models.User, error)
	DeleteUser(id int) error
	GetOwnerID(id int) (int, error)
}

type userService struct {
	UserRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		UserRepo: userRepo,
	}
}

type UserRequest struct {
	Email    string      `json:"email" validate:"required,email"`
	Password string      `json:"password" validate:"required,min=6"`
	Role     models.Role `json:"role" validate:"required,oneof=admin customer"`
}

func (us *userService) GetUserByID(id int) (*models.User, error) {
	return us.UserRepo.GetByID(id)
}

func (us *userService) UpdateUser(id int, req *UserRequest) (*models.User, error) {
	user, err := us.UserRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.Email = req.Email
	user.Role = req.Role
	err = us.UserRepo.Update(user)
	return user, err
}

func (us *userService) UpdateUserPassword(id int, req *UserRequest) (*models.User, error) {
	user, err := us.UserRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err == nil {
		return nil, errors.New("password identical to saved one")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	err = us.UserRepo.Update(user)
	return user, err
}

func (us *userService) DeleteUser(id int) error {
	return us.UserRepo.Delete(id)
}

func (us *userService) GetOwnerID(id int) (int, error) {
	user, err := us.GetUserByID(id)
	if err != nil {
		return 0, err
	}
	return user.ID, err
}
