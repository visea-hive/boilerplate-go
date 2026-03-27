package service

import (
	"github.com/visea-hive/auth-core/internal/model"
	"github.com/visea-hive/auth-core/internal/repository"
	"github.com/visea-hive/auth-core/internal/request"
	"github.com/visea-hive/auth-core/internal/response"
	"github.com/visea-hive/auth-core/pkg/datatable"
	"github.com/visea-hive/auth-core/pkg/messages"
)

type ProductService interface {
	CreateProduct(req request.CreateProductRequest) (response.ProductResponse, error)
	GetProducts() ([]response.ProductResponse, error)
	GetProductsDatatable(dtReq *datatable.Request) ([]response.ProductResponse, int, int, error)
	GetProductByID(id uint) (response.ProductResponse, error)
	UpdateProduct(id uint, req request.UpdateProductRequest) (response.ProductResponse, error)
	DeleteProduct(id uint) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) GetProducts() ([]response.ProductResponse, error) {
	products, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]response.ProductResponse, len(products))
	for i := range products {
		result[i] = response.NewProductResponse(&products[i])
	}
	return result, nil
}

func (s *productService) GetProductsDatatable(dtReq *datatable.Request) ([]response.ProductResponse, int, int, error) {
	products, totalData, totalFiltered, err := s.repo.FindPaginated(dtReq)
	if err != nil {
		return nil, 0, 0, err
	}

	result := make([]response.ProductResponse, len(products))
	for i := range products {
		result[i] = response.NewProductResponse(&products[i])
	}
	return result, totalData, totalFiltered, nil
}


func (s *productService) GetProductByID(id uint) (response.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return response.ProductResponse{}, messages.ErrGeneralNotFound
	}

	return response.NewProductResponse(product), nil
}

func (s *productService) CreateProduct(req request.CreateProductRequest) (response.ProductResponse, error) {
	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.repo.Create(&product); err != nil {
		return response.ProductResponse{}, err
	}

	return response.NewProductResponse(&product), nil
}

func (s *productService) UpdateProduct(id uint, req request.UpdateProductRequest) (response.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return response.ProductResponse{}, messages.ErrGeneralNotFound
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock

	if err := s.repo.Update(product); err != nil {
		return response.ProductResponse{}, err
	}

	return response.NewProductResponse(product), nil
}

func (s *productService) DeleteProduct(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return messages.ErrGeneralNotFound
	}

	return s.repo.Delete(id)
}

