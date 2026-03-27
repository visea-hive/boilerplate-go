package datatable

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/visea-hive/auth-core/pkg/helpers"
	"github.com/visea-hive/auth-core/pkg/messages"
	"gorm.io/gorm"
)

// Request represents common parameters for paginated, sortable, and searchable datatable queries.
type Request struct {
	Page                  int
	Limit                 int
	Filter                []map[string]string
	SearchColumns         []string
	SelectedSearchColumns []string
	Search                string
	SortOrder             string
	SortColumn            string
	StartDateFilter       *time.Time
	EndDateFilter         *time.Time
	Status                *string
}

// NewRequestForTemplate returns a Request preset for fetching all records (no pagination).
func NewRequestForTemplate() *Request {
	return &Request{
		Page:                  1,
		Limit:                 0,
		Filter:                []map[string]string{},
		SearchColumns:         []string{},
		SelectedSearchColumns: []string{},
		Search:                "",
		SortOrder:             "asc",
		SortColumn:            "id",
	}
}

// ParseRequest extracts and validates datatable query parameters from the request context.
// searchColumns: allowed columns for searching.
// sortColumns: allowed columns for sorting.
// sortOrders: allowed sort orders (e.g., "asc", "desc").
func ParseRequest(c *gin.Context, searchColumns, sortColumns, sortOrders []string) (*Request, error) {
	// Parse page
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1 // default first page
	}

	// Parse limit
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 0 {
		limit = 10 // default to 10 records
	}

	// Parse search
	search := c.DefaultQuery("search", "")

	// Parse selected search columns
	selectedSearchColumns := searchColumns
	if searchCols := c.DefaultQuery("selectedSearch", ""); searchCols != "" {
		selectedSearchColumns = strings.Split(searchCols, ",")
	}

	// Validate sort column
	sortCol := c.DefaultQuery("sortCol", sortColumns[0])
	if !helpers.ContainsString(sortColumns, sortCol) {
		return nil, messages.ErrSortColumn
	}

	// Validate sort order
	sortOrder := c.DefaultQuery("sortOrder", sortOrders[0])
	if !helpers.ContainsString(sortOrders, sortOrder) {
		return nil, messages.ErrSortOrder
	}

	// FILTER AREA
	var filter []map[string]string
	if businessPartnerId := c.Query("businessPartnerId"); businessPartnerId != "" {
		filter = append(filter, map[string]string{
			"business_partner_id": businessPartnerId,
		})
	}
	if areaId := c.Query("areaId"); areaId != "" {
		filter = append(filter, map[string]string{
			"area_id": areaId,
		})
	}
	if areaId := c.Query("businessPartnerAreaId"); areaId != "" {
		filter = append(filter, map[string]string{
			"business_partner_area_id": areaId,
		})
	}

	// Parse status filter
	var status *string
	if statusParam := c.Query("status"); statusParam != "" {
		status = &statusParam
	}
	// END FILTER AREA

	return &Request{
		Page:                  page,
		Limit:                 limit,
		Search:                search,
		SearchColumns:         searchColumns,
		SelectedSearchColumns: selectedSearchColumns,
		SortColumn:            sortCol,
		SortOrder:             sortOrder,
		Filter:                filter,
		Status:                status,
	}, nil
}

// Response sends a standard JSON datatable response with pagination metadata.
func Response(c *gin.Context, data interface{}, searchColumns, sortColumns []string, totalData, totalFilteredData, limit int) {
	totalPages := 1
	if limit != 0 {
		totalPages = (totalFilteredData + limit - 1) / limit
	}

	c.JSON(http.StatusOK, gin.H{
		"data":              data,
		"searchColumns":     searchColumns,
		"sortColumns":       sortColumns,
		"totalPage":         totalPages,
		"totalData":         totalData,
		"totalFilteredData": totalFilteredData,
	})
}

// ApplySearch modifies a GORM query to include searching based on the datatable Request.
func (dtReq *Request) ApplySearch(query *gorm.DB) *gorm.DB {
	if dtReq.Search != "" && len(dtReq.SearchColumns) > 0 {
		searchQuery := ""
		var searchArgs []interface{}
		for i, col := range dtReq.SearchColumns {
			if i > 0 {
				searchQuery += " OR "
			}
			searchQuery += col + " ILIKE ?"
			searchArgs = append(searchArgs, "%"+dtReq.Search+"%")
		}
		query = query.Where(searchQuery, searchArgs...)
	}
	return query
}

// ApplyPaginationAndSort modifies a GORM query to include sorting and pagination based on the datatable Request.
func (dtReq *Request) ApplyPaginationAndSort(query *gorm.DB) *gorm.DB {
	if dtReq.SortColumn != "" {
		order := "asc"
		if dtReq.SortOrder == "desc" {
			order = "desc"
		}
		query = query.Order(dtReq.SortColumn + " " + order)
	}

	if dtReq.Limit > 0 {
		offset := (dtReq.Page - 1) * dtReq.Limit
		query = query.Offset(offset).Limit(dtReq.Limit)
	}

	return query
}
