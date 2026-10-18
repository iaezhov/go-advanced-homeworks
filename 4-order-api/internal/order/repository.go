package order

import (
	"4-order-api/pkg/db"
)

type OrderRepository struct {
	Database *db.Db
}

func NewOrderRepository(database *db.Db) *OrderRepository {
	return &OrderRepository{
		Database: database,
	}
}

func (repo *OrderRepository) Create(order *Order) (*Order, error) {
	if err := repo.Database.DB.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (repo *OrderRepository) Get(id uint) (*Order, error) {
	var order Order
	if err := repo.Database.DB.Preload("Items").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (repo *OrderRepository) List(payload OrderListPayload) (OrderListResponse, error) {
	var orders []Order
	var total int64

	if err := repo.Database.DB.Model(&Order{}).Where("user_id = ?", payload.UserID).Count(&total).Error; err != nil {
		return OrderListResponse{}, err
	}

	if err := repo.Database.DB.Preload("Items").
		Where("user_id = ?", payload.UserID).
		Limit(payload.Limit).
		Offset(payload.Offset).
		Find(&orders).Error; err != nil {
		return OrderListResponse{}, err
	}

	return OrderListResponse{
		Items: orders,
		Total: int(total),
	}, nil
}
