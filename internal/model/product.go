package model

// TableName specifies the table name for the Product model.
func (Product) TableName() string {
	return "products"
}

type Product struct {
	BaseModel
	Name        string  `gorm:"type:varchar(255);not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int     `gorm:"type:int;not null;default:0" json:"stock"`
}
