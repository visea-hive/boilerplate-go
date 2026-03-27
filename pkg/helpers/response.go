package helpers

// Response is the standard API response envelope.
type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// PaginationMeta holds pagination metadata returned with list responses.
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalData  int `json:"totalData"`
	TotalPages int `json:"totalPages"`
}

// PaginatedData wraps a list payload with its pagination metadata.
type PaginatedData struct {
	Items any            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

// SuccessResponse creates a successful API response.
func SuccessResponse(message string, data any) Response {
	return Response{Status: true, Message: message, Data: data}
}

// ErrorResponse creates an error API response with an optional data payload.
func ErrorResponse(message string, data ...any) Response {
	var respData any
	if len(data) > 0 {
		respData = data[0]
	}
	return Response{Status: false, Message: message, Data: respData}
}

// NewPaginationMeta builds PaginationMeta and calculates total pages.
func NewPaginationMeta(page, limit, totalData int) PaginationMeta {
	totalPages := 1
	if limit > 0 && totalData > 0 {
		totalPages = (totalData + limit - 1) / limit
	}
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalData:  totalData,
		TotalPages: totalPages,
	}
}

// PaginatedResponse creates a success response containing items + pagination metadata.
func PaginatedResponse(message string, items any, page, limit, totalData int) Response {
	return SuccessResponse(message, PaginatedData{
		Items: items,
		Meta:  NewPaginationMeta(page, limit, totalData),
	})
}
