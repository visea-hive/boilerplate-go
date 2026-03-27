package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visea-hive/auth-core/internal/request"
	"github.com/visea-hive/auth-core/internal/service"
	"github.com/visea-hive/auth-core/pkg/datatable"
	"github.com/visea-hive/auth-core/pkg/helpers"
	"github.com/visea-hive/auth-core/pkg/messages"
)

const moduleName = "Organization"
const pluralModuleName = "Organizations"

// OrganizationHandler handles HTTP requests for organizations.
type OrganizationHandler struct {
	orgService service.OrganizationService
}

// NewOrganizationHandler creates a new instance of OrganizationHandler.
func NewOrganizationHandler(orgService service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{orgService: orgService}
}

func (h *OrganizationHandler) GetAll(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	searchColumns := []string{"name", "slug", "description"}
	sortColumns := []string{"id", "name", "created_at"}
	sortOrders := []string{"asc", "desc"}

	dtReq, err := datatable.ParseRequest(c, searchColumns, sortColumns, sortOrders)
	if err != nil || dtReq == nil {
		switch {
		case errors.Is(err, messages.ErrSortColumn):
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.Translate(lang, messages.ErrSortColumn)})
		case errors.Is(err, messages.ErrSortOrder):
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.Translate(lang, messages.ErrSortOrder)})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": messages.Translate(lang, messages.ErrInternalServer)})
		}
		return
	}

	orgs, totalData, totalFiltered, err := h.orgService.GetAll(dtReq, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	datatable.Response(c, orgs, dtReq.SearchColumns, []string{dtReq.SortColumn}, int(totalData), int(totalFiltered), dtReq.Limit)
}

func (h *OrganizationHandler) GetAllDeleted(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	searchColumns := []string{"name", "slug", "description"}
	sortColumns := []string{"id", "name", "created_at", "deleted_at"}
	sortOrders := []string{"asc", "desc"}

	dtReq, err := datatable.ParseRequest(c, searchColumns, sortColumns, sortOrders)
	if err != nil || dtReq == nil {
		switch {
		case errors.Is(err, messages.ErrSortColumn):
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.Translate(lang, messages.ErrSortColumn)})
		case errors.Is(err, messages.ErrSortOrder):
			c.JSON(http.StatusBadRequest, gin.H{"error": messages.Translate(lang, messages.ErrSortOrder)})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": messages.Translate(lang, messages.ErrInternalServer)})
		}
		return
	}

	orgs, totalData, totalFiltered, err := h.orgService.GetAll(dtReq, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	datatable.Response(c, orgs, dtReq.SearchColumns, []string{dtReq.SortColumn}, int(totalData), int(totalFiltered), dtReq.Limit)
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	var req request.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError := helpers.GenerateErrorValidationResponse(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrGeneralInvalidInput), "errors": validationError.Errors})
		return
	}

	if check := h.orgService.ValidateCreateRequest(req); check != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrGeneralInvalidInput), "errors": check.Errors})
		return
	}

	org, err := h.orgService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": messages.Translate(lang, messages.SuccessCreate, moduleName), "data": org})
}

func (h *OrganizationHandler) GetByID(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrBadRequest), "errors": err.Error()})
		return
	}

	org, err := h.orgService.GetByID(id, false)
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": messages.Translate(lang, messages.ErrOrganizationNotFound), "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.Translate(lang, messages.SuccessGet, moduleName), "data": org})
}

func (h *OrganizationHandler) GetDeletedByID(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrBadRequest), "errors": err.Error()})
		return
	}

	org, err := h.orgService.GetByID(id, true)
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": messages.Translate(lang, messages.ErrOrganizationNotFound), "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.Translate(lang, messages.SuccessGet, moduleName), "data": org})
}

func (h *OrganizationHandler) Update(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrBadRequest), "errors": err.Error()})
		return
	}

	var req request.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError := helpers.GenerateErrorValidationResponse(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrGeneralInvalidInput), "errors": validationError.Errors})
		return
	}

	if check := h.orgService.ValidateUpdateRequest(req); check != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrGeneralInvalidInput), "errors": check.Errors})
		return
	}

	org, err := h.orgService.Update(id, req)
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": messages.Translate(lang, messages.ErrOrganizationNotFound), "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.Translate(lang, messages.SuccessUpdate, moduleName), "data": org})
}

func (h *OrganizationHandler) Delete(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrBadRequest), "errors": err.Error()})
		return
	}

	userInfo := helpers.GetUserInformation(c)
	deletedBy := ""
	if userInfo != nil {
		deletedBy = userInfo.UserUUID
	}

	if err := h.orgService.Delete(id, deletedBy); err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": messages.Translate(lang, messages.ErrOrganizationNotFound), "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.Translate(lang, messages.SuccessDelete, moduleName)})
}

func (h *OrganizationHandler) Restore(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrBadRequest), "errors": err.Error()})
		return
	}

	userInfo := helpers.GetUserInformation(c)
	updatedBy := ""
	if userInfo != nil {
		updatedBy = userInfo.UserUUID
	}

	if err := h.orgService.Restore(id, updatedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.Translate(lang, messages.SuccessUpdate, moduleName+" Restored")})
}

func (h *OrganizationHandler) BulkDelete(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	var req request.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError := helpers.GenerateErrorValidationResponse(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrGeneralInvalidInput), "errors": validationError.Errors})
		return
	}

	userInfo := helpers.GetUserInformation(c)
	deletedBy := ""
	if userInfo != nil {
		deletedBy = userInfo.UserUUID
	}

	if err := h.orgService.BulkDelete(req.IDs, deletedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.Translate(lang, messages.SuccessDelete, pluralModuleName)})
}

func (h *OrganizationHandler) BulkRestore(c *gin.Context) {
	lang := messages.ParseLang(c.GetHeader("Accept-Language"))

	var req request.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError := helpers.GenerateErrorValidationResponse(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.Translate(lang, messages.ErrGeneralInvalidInput), "errors": validationError.Errors})
		return
	}

	userInfo := helpers.GetUserInformation(c)
	updatedBy := ""
	if userInfo != nil {
		updatedBy = userInfo.UserUUID
	}

	if err := h.orgService.BulkRestore(req.IDs, updatedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.Translate(lang, messages.ErrInternalServer), "errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.Translate(lang, messages.SuccessUpdate, pluralModuleName+" Restored")})
}
