package services

import (
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/models"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/repositories"
)

type ProductService interface {
	GetProductByID(id int) (*models.Product, error)
	CreateProduct(req *ProductRequest) (*models.Product, error)
	UpdateProduct(id int, req *ProductRequest) (*models.Product, error)
	DeleteProduct(id int) error
	GetAllProducts() ([]*models.Product, error)
}

type productService struct {
	ProductRepo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{
		ProductRepo: repo,
	}
}

type ProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Category    string  `json:"category" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gt=-1"`
}

func (ps *productService) GetProductByID(id int) (*models.Product, error) {
	return ps.ProductRepo.GetByID(id)
}

func (ps *productService) CreateProduct(req *ProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	err := ps.ProductRepo.Create(product)
	return product, err
}

func (ps *productService) UpdateProduct(id int, req *ProductRequest) (*models.Product, error) {
	product, err := ps.ProductRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	product.Name = req.Name
	product.Description = req.Description
	product.Category = req.Category
	product.Price = req.Price
	product.Stock = req.Stock
	err = ps.ProductRepo.Update(product)
	return product, err
}

func (ps *productService) DeleteProduct(id int) error {
	return ps.ProductRepo.Delete(id)
}

func (ps *productService) GetAllProducts() ([]*models.Product, error) {
	return ps.ProductRepo.GetAll()
}
