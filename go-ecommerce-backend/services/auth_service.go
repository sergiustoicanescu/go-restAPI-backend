package services

import (
	"errors"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/repositories"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req *LoginRequest) (*models.User, string, error)
	Register(req *RegisterRequest) error
}

type authService struct {
	UserRepo  repositories.UserRepository
	JWTSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (a *authService) Login(req *LoginRequest) (*models.User, string, error) {
	user, err := a.UserRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, string(user.Role), a.JWTSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

type RegisterRequest struct {
	Email    string      `json:"email" validate:"required,email"`
	Password string      `json:"password" validate:"required,min=6"`
	Role     models.Role `json:"role" validate:"required"`
}

func (a *authService) Register(req *RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}
	return a.UserRepo.Create(user)
}
