package repositories

import (
	"database/sql"
	"fmt"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
)

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id int) (*models.Order, error)
	GetOrdersByCustomerID(id int) ([]*models.Order, error)
	Update(order *models.Order) error
	GetOwnerID(id int) (int, error)
}

type orderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{DB: db}
}

func (or *orderRepository) Create(order *models.Order) error {
	tx, err := or.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	orderQuery := "INSERT INTO orders (customer_id, status) VALUES ($1, $2) RETURNING id, created_at"

	err = tx.QueryRow(orderQuery, order.CustomerID, order.Status).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return err
	}

	for _, item := range order.OrderItems {
		var available int
		stockQuery := "SELECT stock FROM products WHERE id = $1 FOR UPDATE"
		err = tx.QueryRow(stockQuery, item.ProductID).Scan(&available)
		if err != nil {
			return err
		}

		if available < item.Quantity {
			return fmt.Errorf("insufficient quantity for product %d: available %d, required %d", item.ProductID, available, item.Quantity)
		}

		updateStockQuery := "UPDATE products SET stock = stock - $1 WHERE id = $2"
		_, err = tx.Exec(updateStockQuery, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}

		itemQuery := "INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)"
		_, err = tx.Exec(itemQuery, order.ID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (or *orderRepository) GetByID(id int) (*models.Order, error) {
	order := &models.Order{}
	query := "SELECT id, customer_id, status, created_at FROM orders WHERE id = $1"
	err := or.DB.QueryRow(query, id).Scan(&order.ID, &order.CustomerID, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	itemsQuery := "SELECT id, product_id, quantity FROM order_items WHERE order_id = $1"

	rows, err := or.DB.Query(itemsQuery, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err = rows.Scan(&item.ID, &item.ProductID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		item.OrderID = order.ID
		orderItems = append(orderItems, item)
	}
	order.OrderItems = orderItems

	return order, nil
}

func (or *orderRepository) GetOrdersByCustomerID(id int) ([]*models.Order, error) {
	query := "SELECT id, customer_id, status, created_at FROM orders WHERE customer_id = $1"
	rows, err := or.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order = new(models.Order)
		if err := rows.Scan(&order.ID, &order.CustomerID, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (or *orderRepository) Update(order *models.Order) error {
	query := "UPDATE orders SET status = $1 WHERE id = $2"
	_, err := or.DB.Exec(query, order.Status, order.ID)
	return err
}

func (or *orderRepository) GetOwnerID(id int) (int, error) {
	var userID int
	query := "SELECT c.user_id FROM orders o JOIN customers c ON c.id = o.customer_id WHERE o.id = $1"
	err := or.DB.QueryRow(query, id).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, err
}
