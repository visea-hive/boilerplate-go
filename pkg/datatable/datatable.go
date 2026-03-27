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
//
//   - searchColumns: allowed columns for searching
//   - sortColumns:   allowed columns for sorting
//   - sortOrders:    allowed sort orders (e.g. "asc", "desc")
//   - filterKeys:    optional extra query-param keys to capture as column filters
//     (camelCase keys are kept as-is; map them to DB columns in the repository layer)
func ParseRequest(c *gin.Context, searchColumns, sortColumns, sortOrders []string, filterKeys ...string) (*Request, error) {
	// Parse page
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	// Parse limit
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 0 {
		limit = 10
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

	// Collect caller-supplied filter keys from query params
	var filter []map[string]string
	for _, key := range filterKeys {
		if val := c.Query(key); val != "" {
			filter = append(filter, map[string]string{key: val})
		}
	}

	return &Request{
		Page:                  page,
		Limit:                 limit,
		Search:                search,
		SearchColumns:         searchColumns,
		SelectedSearchColumns: selectedSearchColumns,
		SortColumn:            sortCol,
		SortOrder:             sortOrder,
		Filter:                filter,
	}, nil
}

// ApplySearch modifies a GORM query to include ILIKE search across SearchColumns.
func (dtReq *Request) ApplySearch(query *gorm.DB) *gorm.DB {
	if dtReq.Search != "" && len(dtReq.SearchColumns) > 0 {
		var parts []string
		var args []any
		for _, col := range dtReq.SearchColumns {
			parts = append(parts, col+" ILIKE ?")
			args = append(args, "%"+dtReq.Search+"%")
		}
		query = query.Where(strings.Join(parts, " OR "), args...)
	}
	return query
}

// ApplyFilters applies the key/value pairs in Filter as exact-match WHERE clauses.
func (dtReq *Request) ApplyFilters(query *gorm.DB) *gorm.DB {
	for _, f := range dtReq.Filter {
		for col, val := range f {
			query = query.Where(col+" = ?", val)
		}
	}
	return query
}

// ApplyPaginationAndSort modifies a GORM query to include sorting and pagination.
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

// Response sends a standard JSON datatable response with pagination metadata.
func Response(c *gin.Context, data any, searchColumns, sortColumns []string, totalData, totalFilteredData, limit int) {
	totalPages := 1
	if limit > 0 && totalFilteredData > 0 {
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
