package services

import (
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/repositories"
)

type OrderService interface {
	GetOrderByID(id int) (*models.Order, error)
	CreateOrder(req *OrderRequest) (*models.Order, error)
	UpdateOrder(id int, req *OrderRequest) (*models.Order, error)
	GetOwnerID(id int) (int, error)
}

type orderService struct {
	OrderRepo repositories.OrderRepository
}

func NewOrderService(repo repositories.OrderRepository) OrderService {
	return &orderService{
		OrderRepo: repo,
	}
}

type OrderRequest struct {
	CustomerID int                `json:"customer_id" validate:"required"`
	Status     models.OrderStatus `json:"status" validate:"required,oneof=pending completed cancelled"`
	OrderItems []OrderItemRequest `json:"order_items" validate:"required"`
}

type OrderItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
}

func (os *orderService) GetOrderByID(id int) (*models.Order, error) {
	return os.OrderRepo.GetByID(id)
}

func (os *orderService) CreateOrder(req *OrderRequest) (*models.Order, error) {
	var orderItems []models.OrderItem
	for _, item := range req.OrderItems {
		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}
	order := &models.Order{
		CustomerID: req.CustomerID,
		Status:     req.Status,
		OrderItems: orderItems,
	}
	err := os.OrderRepo.Create(order)
	return order, err
}

func (os *orderService) UpdateOrder(id int, req *OrderRequest) (*models.Order, error) {
	order, err := os.OrderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	order.Status = req.Status
	err = os.OrderRepo.Update(order)
	return order, err
}

func (os *orderService) GetOwnerID(id int) (int, error) {
	return os.OrderRepo.GetOwnerID(id)
}
