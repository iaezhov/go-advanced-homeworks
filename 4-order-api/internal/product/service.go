package product

import (
	"4-order-api/pkg/db"
)

type ProductService struct {
	db *db.Db
}

func (s *ProductService) GetByID(id uint) (*Product, error) {
	var p Product
	result := s.db.First(&p, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &p, nil
}

func NewProductService(db *db.Db) *ProductService {
	return &ProductService{db: db}
}
