package repository

import (
	"github.com/visea-hive/auth-core/internal/model"
	"github.com/visea-hive/auth-core/pkg/datatable"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *model.Product) error
	FindAll() ([]model.Product, error)
	FindByID(id uint) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id uint) error
	FindPaginated(dtReq *datatable.Request) (products []model.Product, totalData int, totalFiltered int, err error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindPaginated(dtReq *datatable.Request) ([]model.Product, int, int, error) {
	var products []model.Product
	var totalData, totalFiltered int64

	// Count all records before any filter
	if err := r.db.Model(&model.Product{}).Count(&totalData).Error; err != nil {
		return nil, 0, 0, err
	}

	// Build search + filter query
	query := r.db.Model(&model.Product{})
	query = dtReq.ApplySearch(query)
	query = dtReq.ApplyFilters(query)

	// Count after search/filters
	if err := query.Count(&totalFiltered).Error; err != nil {
		return nil, 0, 0, err
	}

	// Fetch paginated, sorted results
	if err := dtReq.ApplyPaginationAndSort(query).Find(&products).Error; err != nil {
		return nil, 0, 0, err
	}

	return products, int(totalData), int(totalFiltered), nil
}


func (r *productRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}