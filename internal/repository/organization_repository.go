package repository

import (
	"time"

	"github.com/visea-hive/auth-core/internal/model"
	"github.com/visea-hive/auth-core/pkg/datatable"
	"gorm.io/gorm"
)

// OrganizationRepository defines the interface for organization-related database operations.
type OrganizationRepository interface {
	GetAll(db *gorm.DB, dtReq *datatable.Request, fetchDeleted bool) ([]model.OrganizationDTO, int64, int64, error)
	Create(db *gorm.DB, org *model.Organization) error
	GetByID(db *gorm.DB, id uint) (*model.Organization, error)
	GetDeletedByID(db *gorm.DB, id uint) (*model.Organization, error)
	GetBySlug(db *gorm.DB, slug string) (*model.Organization, error)
	Update(db *gorm.DB, org *model.Organization) error
	Delete(db *gorm.DB, id uint, deletedBy string) error
	Restore(db *gorm.DB, id uint, updatedBy string) error
	BulkDelete(db *gorm.DB, ids []uint, deletedBy string) error
	BulkRestore(db *gorm.DB, ids []uint, updatedBy string) error
}

type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new instance of OrganizationRepository.
func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) GetAll(db *gorm.DB, dtReq *datatable.Request, fetchDeleted bool) ([]model.OrganizationDTO, int64, int64, error) {
	query := db.Model(&model.Organization{})

	if fetchDeleted {
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	} else {
		query = query.Where("status != ?", model.StatusSuspended)
	}

	// 1. Base count without datatable text filters
	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, 0, err
	}

	// 2. Apply Datatable UI Search
	query = dtReq.ApplySearch(query)

	// -- You can apply custom repository filters using query.Where() here --

	// 3. Count data after all filters are applied
	var totalFiltered int64
	if err := query.Count(&totalFiltered).Error; err != nil {
		return nil, 0, 0, err
	}

	// 4. Apply Datatable Pagination & Sorting before finding
	query = dtReq.ApplyPaginationAndSort(query)

	var orgs []model.Organization
	if err := query.Preload("Parent").Find(&orgs).Error; err != nil {
		return nil, 0, 0, err
	}

	dtos := model.MapToOrganizationDTOs(orgs)

	return dtos, totalData, totalFiltered, nil
}

func (r *organizationRepository) Create(db *gorm.DB, org *model.Organization) error {
	return db.Create(org).Error
}

func (r *organizationRepository) GetByID(db *gorm.DB, id uint) (*model.Organization, error) {
	var org model.Organization
	if err := db.Preload("Parent").First(&org, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return explicit nil on not found
		}
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) GetDeletedByID(db *gorm.DB, id uint) (*model.Organization, error) {
	var org model.Organization
	if err := db.Unscoped().Preload("Parent").Where("deleted_at IS NOT NULL").First(&org, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return explicit nil on not found
		}
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) GetBySlug(db *gorm.DB, slug string) (*model.Organization, error) {
	var org model.Organization
	if err := db.Preload("Parent").Where("slug = ?", slug).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return explicit nil on not found
		}
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) Update(db *gorm.DB, org *model.Organization) error {
	return db.Save(org).Error
}

func (r *organizationRepository) Delete(db *gorm.DB, id uint, deletedBy string) error {
	return db.Model(&model.Organization{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_by": deletedBy,
		"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
	}).Error
}

func (r *organizationRepository) Restore(db *gorm.DB, id uint, updatedBy string) error {
	return db.Unscoped().Model(&model.Organization{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_by": nil,
		"deleted_at": gorm.DeletedAt{Valid: false},
		"updated_by": updatedBy,
	}).Error
}

func (r *organizationRepository) BulkDelete(db *gorm.DB, ids []uint, deletedBy string) error {
	return db.Model(&model.Organization{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"deleted_by": deletedBy,
		"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
	}).Error
}

func (r *organizationRepository) BulkRestore(db *gorm.DB, ids []uint, updatedBy string) error {
	return db.Unscoped().Model(&model.Organization{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"deleted_by": nil,
		"deleted_at": gorm.DeletedAt{Valid: false},
		"updated_by": updatedBy,
	}).Error
}
