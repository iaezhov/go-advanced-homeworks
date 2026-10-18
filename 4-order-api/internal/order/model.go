package order

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID uint    `json:"user_id"`
	Price  float64 `json:"price"`
	Items  []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
