package product

import "github.com/lib/pq"

type ProductCreatePayload struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Price       float64        `json:"price" validate:"required"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
}
type ProductUpdatePayload struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
}
