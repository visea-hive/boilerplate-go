package service

import (
	"github.com/visea-hive/auth-core/internal/model"
	"github.com/visea-hive/auth-core/internal/repository"
	"github.com/visea-hive/auth-core/internal/request"
	"github.com/visea-hive/auth-core/pkg/datatable"
	"github.com/visea-hive/auth-core/pkg/helpers"
	"github.com/visea-hive/auth-core/pkg/messages"
	"gorm.io/gorm"
)

// OrganizationService defines the interface for organization business logic.
type OrganizationService interface {
	GetAll(dtReq *datatable.Request, fetchDeleted bool) ([]model.OrganizationDTO, int64, int64, error)
	GetByID(id uint, fetchDeleted bool) (model.OrganizationDTO, error)
	Create(req request.CreateOrganizationRequest) (model.OrganizationDTO, error)
	Update(id uint, req request.UpdateOrganizationRequest) (model.OrganizationDTO, error)
	Delete(id uint, deletedBy string) error
	Restore(id uint, updatedBy string) error
	BulkDelete(ids []uint, deletedBy string) error
	BulkRestore(ids []uint, updatedBy string) error
	ValidateCreateRequest(req request.CreateOrganizationRequest) *helpers.ErrorResponse
	ValidateUpdateRequest(req request.UpdateOrganizationRequest) *helpers.ErrorResponse
}

type organizationService struct {
	orgRepo repository.OrganizationRepository
	db      *gorm.DB
}

// NewOrganizationService creates a new instance of OrganizationService.
func NewOrganizationService(orgRepo repository.OrganizationRepository, db *gorm.DB) OrganizationService {
	return &organizationService{
		orgRepo: orgRepo,
		db:      db,
	}
}

func (s *organizationService) GetAll(dtReq *datatable.Request, fetchDeleted bool) ([]model.OrganizationDTO, int64, int64, error) {
	return s.orgRepo.GetAll(s.db, dtReq, fetchDeleted)
}

func (s *organizationService) Create(req request.CreateOrganizationRequest) (model.OrganizationDTO, error) {
	// Generate base slug from name
	slug := helpers.GenerateSlug(req.Name)

	// In a real app, you'd check for slug collision and append numbers if necessary.
	existing, err := s.orgRepo.GetBySlug(s.db, slug)
	if err != nil {
		return model.OrganizationDTO{}, err
	}
	if existing != nil {
		return model.OrganizationDTO{}, messages.ErrGeneralCreated
	}

	org := &model.Organization{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Level:       req.Level,
		ParentID:    req.ParentID,
		Status:      model.StatusActive,
		SSOEnforced: false,
		MFAEnforced: false,
	}

	err = s.db.Transaction(func(db *gorm.DB) error {
		return s.orgRepo.Create(db, org)
	})
	if err != nil {
		return model.OrganizationDTO{}, messages.ErrGeneralCreated
	}

	return model.MapToOrganizationDTO(org), nil
}

func (s *organizationService) GetByID(id uint, fetchDeleted bool) (model.OrganizationDTO, error) {
	var org *model.Organization
	var err error

	if fetchDeleted {
		org, err = s.orgRepo.GetDeletedByID(s.db, id)
	} else {
		org, err = s.orgRepo.GetByID(s.db, id)
	}

	if err != nil {
		return model.OrganizationDTO{}, err
	}
	if org == nil {
		return model.OrganizationDTO{}, messages.ErrOrganizationNotFound
	}

	return model.MapToOrganizationDTO(org), nil
}

func (s *organizationService) Update(id uint, req request.UpdateOrganizationRequest) (model.OrganizationDTO, error) {
	org, err := s.orgRepo.GetByID(s.db, id)
	if err != nil {
		return model.OrganizationDTO{}, err
	}
	if org == nil {
		return model.OrganizationDTO{}, messages.ErrOrganizationNotFound
	}

	// Update fields if provided
	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Description != nil {
		org.Description = req.Description
	}
	if req.Status != nil {
		org.Status = model.OrganizationStatus(*req.Status)
	}
	if req.SSOEnforced != nil {
		org.SSOEnforced = *req.SSOEnforced
	}
	if req.MFAEnforced != nil {
		org.MFAEnforced = *req.MFAEnforced
	}

	err = s.db.Transaction(func(db *gorm.DB) error {
		return s.orgRepo.Update(db, org)
	})
	if err != nil {
		return model.OrganizationDTO{}, messages.ErrGeneralUpdated
	}

	return model.MapToOrganizationDTO(org), nil
}

func (s *organizationService) Delete(id uint, deletedBy string) error {
	// Let's first ensure it exists
	org, err := s.orgRepo.GetByID(s.db, id)
	if err != nil {
		return err
	}
	if org == nil {
		return messages.ErrOrganizationNotFound
	}

	err = s.db.Transaction(func(db *gorm.DB) error {
		return s.orgRepo.Delete(db, id, deletedBy)
	})
	if err != nil {
		return messages.ErrGeneralDeleted
	}

	return nil
}

func (s *organizationService) Restore(id uint, updatedBy string) error {
	org, err := s.orgRepo.GetDeletedByID(s.db, id)
	if err != nil {
		return err
	}
	if org == nil {
		return messages.ErrOrganizationNotFound
	}

	err = s.db.Transaction(func(db *gorm.DB) error {
		return s.orgRepo.Restore(db, id, updatedBy)
	})
	if err != nil {
		return messages.ErrGeneralUpdated
	}

	return nil
}

func (s *organizationService) BulkDelete(ids []uint, deletedBy string) error {
	err := s.db.Transaction(func(db *gorm.DB) error {
		return s.orgRepo.BulkDelete(db, ids, deletedBy)
	})
	if err != nil {
		return messages.ErrGeneralDeleted
	}
	return nil
}

func (s *organizationService) BulkRestore(ids []uint, updatedBy string) error {
	err := s.db.Transaction(func(db *gorm.DB) error {
		return s.orgRepo.BulkRestore(db, ids, updatedBy)
	})
	if err != nil {
		return messages.ErrGeneralUpdated
	}
	return nil
}

func (s *organizationService) ValidateCreateRequest(req request.CreateOrganizationRequest) *helpers.ErrorResponse {
	response := helpers.ErrorResponse{}
	existing, err := s.orgRepo.GetBySlug(s.db, helpers.GenerateSlug(req.Name))
	if err == nil && existing != nil {
		response.Errors = append(response.Errors, helpers.ValidationError{
			FieldName:   "name",
			RuleMessage: messages.ErrOrganizationNameExists.Error(),
		})
	}
	if len(response.Errors) > 0 {
		return &response
	}
	return nil
}

func (s *organizationService) ValidateUpdateRequest(req request.UpdateOrganizationRequest) *helpers.ErrorResponse {
	response := helpers.ErrorResponse{}
	if req.Name != nil {
		existing, err := s.orgRepo.GetBySlug(s.db, helpers.GenerateSlug(*req.Name))
		if err == nil && existing != nil && existing.Name == *req.Name {
			response.Errors = append(response.Errors, helpers.ValidationError{
				FieldName:   "name",
				RuleMessage: messages.ErrOrganizationNameExists.Error(),
			})
		}
	}
	if len(response.Errors) > 0 {
		return &response
	}
	return nil
}
