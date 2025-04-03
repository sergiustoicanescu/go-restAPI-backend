package repositories

import (
	"database/sql"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id int) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id int) error
	GetAll() ([]*models.Product, error)
}

type productRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{DB: db}
}

func (pr *productRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, description, category, price, stock) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	return pr.DB.QueryRow(query, product.Name, product.Description, product.Category, product.Price, product.Stock).Scan(&product.ID)
}

func (pr *productRepository) GetByID(id int) (*models.Product, error) {
	product := &models.Product{}
	query := "SELECT id, name, description, category, price, stock FROM products WHERE id = $1"
	err := pr.DB.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Stock)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pr *productRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, description = $2, category = $3, price = $4, stock = $5 WHERE id = $6"
	_, err := pr.DB.Exec(query, product.Name, product.Description, product.Category, product.Price, product.Stock, product.ID)
	return err
}

func (pr *productRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := pr.DB.Exec(query, id)
	return err
}

func (pr *productRepository) GetAll() ([]*models.Product, error) {
	query := "SELECT id, name, description, category, price, stock FROM products"
	rows, err := pr.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product = new(models.Product)
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
