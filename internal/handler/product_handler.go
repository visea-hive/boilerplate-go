package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/visea-hive/auth-core/internal/middleware"
	"github.com/visea-hive/auth-core/internal/request"
	"github.com/visea-hive/auth-core/internal/response"
	"github.com/visea-hive/auth-core/internal/service"
	"github.com/visea-hive/auth-core/pkg/datatable"
	"github.com/visea-hive/auth-core/pkg/helpers"
	"github.com/visea-hive/auth-core/pkg/messages"
)

// Allowed columns for product datatable queries.
var (
	productSearchColumns = []string{"name", "description"}
	productSortColumns   = []string{"id", "name", "price", "stock", "created_at"}
	productSortOrders    = []string{"asc", "desc"}
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

func (h *ProductHandler) Create(c *gin.Context) {
	lang := middleware.GetLang(c)

	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.GenerateErrorValidationResponse(err))
		return
	}

	product, err := h.service.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrorResponse(messages.Translate(lang, err)))
		return
	}

	c.JSON(http.StatusCreated, helpers.SuccessResponse(messages.Translate(lang, messages.SuccessCreate, "product"), product))
}

func (h *ProductHandler) FindAll(c *gin.Context) {
	lang := middleware.GetLang(c)

	products, err := h.service.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrorResponse(messages.Translate(lang, err)))
		return
	}

	if products == nil {
		products = []response.ProductResponse{}
	}

	c.JSON(http.StatusOK, helpers.SuccessResponse(messages.Translate(lang, messages.SuccessGet, "product"), products))
}

// FindPaginated handles GET /products/datatable
// Supports: ?page=1&limit=10&search=foo&sortCol=name&sortOrder=asc&selectedSearch=name,description
func (h *ProductHandler) FindPaginated(c *gin.Context) {
	lang := middleware.GetLang(c)

	dtReq, err := datatable.ParseRequest(c, productSearchColumns, productSortColumns, productSortOrders)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(messages.Translate(lang, err)))
		return
	}

	products, totalData, totalFiltered, err := h.service.GetProductsDatatable(dtReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrorResponse(messages.Translate(lang, err)))
		return
	}

	datatable.Response(c, products, productSearchColumns, productSortColumns, totalData, totalFiltered, dtReq.Limit)
}

func (h *ProductHandler) FindByID(c *gin.Context) {
	lang := middleware.GetLang(c)

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(messages.Translate(lang, messages.ErrGeneralID)))
		return
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrorResponse(messages.Translate(lang, messages.ErrGeneralNotFound)))
		return
	}

	c.JSON(http.StatusOK, helpers.SuccessResponse(messages.Translate(lang, messages.SuccessGet, "product"), product))
}

func (h *ProductHandler) Update(c *gin.Context) {
	lang := middleware.GetLang(c)

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(messages.Translate(lang, messages.ErrGeneralID)))
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.GenerateErrorValidationResponse(err))
		return
	}

	product, err := h.service.UpdateProduct(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrorResponse(messages.Translate(lang, err)))
		return
	}

	c.JSON(http.StatusOK, helpers.SuccessResponse(messages.Translate(lang, messages.SuccessUpdate, "product"), product))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	lang := middleware.GetLang(c)

	id, err := helpers.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(messages.Translate(lang, messages.ErrGeneralID)))
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrorResponse(messages.Translate(lang, err)))
		return
	}

	c.JSON(http.StatusOK, helpers.SuccessResponse(messages.Translate(lang, messages.SuccessDelete, "product"), nil))
}
