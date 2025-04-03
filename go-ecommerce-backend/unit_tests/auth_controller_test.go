package unit_tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/controllers"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/services"
)

type MockAuthService struct {
	LoginFunc    func(req *services.LoginRequest) (*models.User, string, error)
	RegisterFunc func(req *services.RegisterRequest) error
}

func (m *MockAuthService) Login(req *services.LoginRequest) (*models.User, string, error) {
	return m.LoginFunc(req)
}

func (m *MockAuthService) Register(req *services.RegisterRequest) error {
	return m.RegisterFunc(req)
}

func TestAuthController_Login_Success(t *testing.T) {
	mockUser := &models.User{
		ID:    1,
		Email: "test@example.com",
		Role:  models.RoleCustomer,
	}
	expectedToken := "fake-jwt-token"

	mockService := &MockAuthService{
		LoginFunc: func(req *services.LoginRequest) (*models.User, string, error) {
			return mockUser, expectedToken, nil
		},
	}

	authController := controllers.NewAuthController(mockService)

	loginReq := &services.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	rr := httptest.NewRecorder()

	authController.Login(rr, nil, loginReq)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200 on successful login")
	var respBody map[string]any
	err := json.NewDecoder(rr.Body).Decode(&respBody)
	assert.NoError(t, err, "Expected valid JSON response")

	userMap, ok := respBody["user"].(map[string]any)
	assert.True(t, ok, "User should be a JSON object")
	assert.Equal(t, float64(mockUser.ID), userMap["id"], "User ID should match")
	assert.Equal(t, mockUser.Email, userMap["email"], "User email should match")
	assert.Equal(t, string(mockUser.Role), userMap["role"], "User role should match")
	assert.Equal(t, expectedToken, respBody["token"], "Token should match")
}

func TestAuthController_Login_Failure(t *testing.T) {
	mockService := &MockAuthService{
		LoginFunc: func(req *services.LoginRequest) (*models.User, string, error) {
			return nil, "", errors.New("invalid credentials")
		},
	}

	authController := controllers.NewAuthController(mockService)

	loginReq := &services.LoginRequest{
		Email:    "wrong@example.com",
		Password: "wrongpassword",
	}

	rr := httptest.NewRecorder()
	authController.Login(rr, nil, loginReq)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status code 401 on failed login")
}

func TestAuthController_Register_Success(t *testing.T) {
	mockService := &MockAuthService{
		RegisterFunc: func(req *services.RegisterRequest) error {
			return nil
		},
	}

	authController := controllers.NewAuthController(mockService)

	registerReq := &services.RegisterRequest{
		Email:    "new@example.com",
		Password: "password",
		Role:     models.RoleCustomer,
	}

	rr := httptest.NewRecorder()
	authController.Register(rr, nil, registerReq)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code 201 on successful registration")

	var respBody map[string]string
	err := json.NewDecoder(rr.Body).Decode(&respBody)
	assert.NoError(t, err, "Expected valid JSON response")
	assert.Equal(t, "User registered successfully", respBody["message"], "Expected success message")
}

func TestAuthController_Register_Failure(t *testing.T) {
	mockService := &MockAuthService{
		RegisterFunc: func(req *services.RegisterRequest) error {
			return errors.New("registration error")
		},
	}

	authController := controllers.NewAuthController(mockService)

	registerReq := &services.RegisterRequest{
		Email:    "fail@example.com",
		Password: "password",
		Role:     models.RoleCustomer,
	}

	rr := httptest.NewRecorder()
	authController.Register(rr, nil, registerReq)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status code 500 on registration failure")
}
