package unit_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/controllers"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/services"
)

type MockOrderService struct {
	GetOrderByIDFunc          func(id int) (*models.Order, error)
	GetOrdersByCustomerIDFunc func(id int) ([]*models.Order, error)
	CreateOrderFunc           func(req *services.OrderRequest) (*models.Order, error)
	UpdateOrderFunc           func(id int, req *services.OrderRequest) (*models.Order, error)
	GetOwnerIDFunc            func(id int) (int, error)
}

func (m *MockOrderService) GetOrderByID(id int) (*models.Order, error) {
	return m.GetOrderByIDFunc(id)
}

func (m *MockOrderService) GetOrdersByCustomerID(id int) ([]*models.Order, error) {
	return m.GetOrdersByCustomerIDFunc(id)
}

func (m *MockOrderService) CreateOrder(req *services.OrderRequest) (*models.Order, error) {
	return m.CreateOrderFunc(req)
}

func (m *MockOrderService) UpdateOrder(id int, req *services.OrderRequest) (*models.Order, error) {
	return m.UpdateOrderFunc(id, req)
}

func (m *MockOrderService) GetOwnerID(id int) (int, error) {
	return m.GetOwnerIDFunc(id)
}

var expectedOrder = models.Order{
	ID:         1,
	CustomerID: 1,
	Status:     models.OrderStatusPending,
	OrderItems: []models.OrderItem{
		{
			ID:        1,
			OrderID:   1,
			ProductID: 1,
			Quantity:  2,
		},
		{
			ID:        2,
			OrderID:   1,
			ProductID: 2,
			Quantity:  1,
		},
	},
	CreatedAt: time.Now(),
}

func TestOrderController_GetOrder_Success(t *testing.T) {
	mockService := &MockOrderService{
		GetOrderByIDFunc: func(id int) (*models.Order, error) {
			return &expectedOrder, nil
		},
	}
	orderController := controllers.NewOrderController(mockService)

	req := httptest.NewRequest("GET", "/orders/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	orderController.GetOrder(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200 for a valid order")
	var respOrder models.Order
	err := json.NewDecoder(rr.Body).Decode(&respOrder)
	assert.NoError(t, err, "Expected valid JSON response")
	assert.Equal(t, expectedOrder.ID, respOrder.ID, "Order ID should match")
	assert.Equal(t, expectedOrder.CustomerID, respOrder.CustomerID, "Customer ID should match")
	assert.Equal(t, expectedOrder.Status, respOrder.Status, "Order status should match")
	assert.Equal(t, len(expectedOrder.OrderItems), len(respOrder.OrderItems), "Order items count should match")
	assert.Equal(t, expectedOrder.OrderItems[0].ProductID, respOrder.OrderItems[0].ProductID, "First order item product ID should match")
}

func TestOrderController_GetOrder_InvalidID(t *testing.T) {
	mockService := &MockOrderService{
		GetOrderByIDFunc: func(id int) (*models.Order, error) {
			return &models.Order{}, nil
		},
	}
	orderController := controllers.NewOrderController(mockService)
	req := httptest.NewRequest("GET", "/orders/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	orderController.GetOrder(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code 400 for invalid order ID")
}

func TestOrderController_GetOrder_NotFound(t *testing.T) {
	mockService := &MockOrderService{
		GetOrderByIDFunc: func(id int) (*models.Order, error) {
			return &models.Order{}, errors.New("order not found")
		},
	}
	orderController := controllers.NewOrderController(mockService)
	req := httptest.NewRequest("GET", "/orders/2", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "2"})
	rr := httptest.NewRecorder()

	orderController.GetOrder(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "Expected status code 404 for order not found")
}

func TestOrderController_CreateOrder_Success(t *testing.T) {
	mockService := &MockOrderService{
		CreateOrderFunc: func(req *services.OrderRequest) (*models.Order, error) {
			return &expectedOrder, nil
		},
	}
	orderController := controllers.NewOrderController(mockService)

	orderReq := services.OrderRequest{
		CustomerID: 1,
		Status:     models.OrderStatusPending,
		OrderItems: []services.OrderItemRequest{
			{
				ProductID: 1,
				Quantity:  2,
			},
			{
				ProductID: 2,
				Quantity:  1,
			},
		},
	}
	body, err := json.Marshal(orderReq)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	orderController.CreateOrder(rr, req, &orderReq)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code 201 on successful order creation")
	var respOrder models.Order
	err = json.NewDecoder(rr.Body).Decode(&respOrder)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.ID, respOrder.ID, "Order ID should match")
}

func TestOrderController_UpdateOrder_Success(t *testing.T) {
	localExpectedOrder := expectedOrder
	localExpectedOrder.Status = models.OrderStatusCancelled

	mockService := &MockOrderService{
		UpdateOrderFunc: func(id int, req *services.OrderRequest) (*models.Order, error) {
			return &localExpectedOrder, nil
		},
	}
	orderController := controllers.NewOrderController(mockService)

	orderReq := services.OrderRequest{
		CustomerID: 1,
		Status:     models.OrderStatusCancelled,
		OrderItems: []services.OrderItemRequest{
			{
				ProductID: 1,
				Quantity:  2,
			},
			{
				ProductID: 2,
				Quantity:  1,
			},
		},
	}
	body, err := json.Marshal(orderReq)
	assert.NoError(t, err)

	req := httptest.NewRequest("PUT", "/orders/1", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	orderController.UpdateOrder(rr, req, &orderReq)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200 on successful order update")
	var respOrder models.Order
	err = json.NewDecoder(rr.Body).Decode(&respOrder)
	assert.NoError(t, err)
	assert.Equal(t, localExpectedOrder.ID, respOrder.ID, "Order ID should match")
	assert.Equal(t, localExpectedOrder.Status, respOrder.Status, "Order status should match")
}

func TestOrderController_GetOrdersByCustomerID_Success(t *testing.T) {
	order1 := &models.Order{
		ID:         1,
		CustomerID: 1,
		Status:     models.OrderStatusPending,
		OrderItems: []models.OrderItem{
			{ID: 1, OrderID: 1, ProductID: 1, Quantity: 2},
		},
		CreatedAt: time.Now(),
	}
	order2 := &models.Order{
		ID:         2,
		CustomerID: 1,
		Status:     models.OrderStatusPending,
		OrderItems: []models.OrderItem{
			{ID: 2, OrderID: 2, ProductID: 2, Quantity: 1},
		},
		CreatedAt: time.Now(),
	}
	expectedOrders := []*models.Order{order1, order2}

	mockService := &MockOrderService{
		GetOrdersByCustomerIDFunc: func(id int) ([]*models.Order, error) {
			if id == 1 {
				return expectedOrders, nil
			}
			return nil, errors.New("customer not found")
		},
	}

	orderController := controllers.NewOrderController(mockService)

	req := httptest.NewRequest("GET", "/customers/1/orders", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	orderController.GetOrdersByCustomerID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200 for successful orders retrieval")

	var respOrders []*models.Order
	err := json.NewDecoder(rr.Body).Decode(&respOrders)
	assert.NoError(t, err, "Expected valid JSON response")
	assert.Equal(t, len(expectedOrders), len(respOrders), "Expected orders count to match")
	assert.Equal(t, expectedOrders[0].ID, respOrders[0].ID, "Order ID should match for the first order")
	assert.Equal(t, expectedOrders[1].ID, respOrders[1].ID, "Order ID should match for the second order")
}
