package request

type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
}

type CreateProductRequest struct {
	ProductRequest
}

type UpdateProductRequest struct {
	ProductRequest
}
