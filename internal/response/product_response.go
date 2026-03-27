package response

import (
	"time"

	"github.com/visea-hive/auth-core/internal/model"
)

type ProductResponse struct {
	ID          uint       `json:"id"`
	UUID        string     `json:"uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// NewProductResponse maps a Product model to its API response representation.
func NewProductResponse(p *model.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		UUID:        p.UUID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
