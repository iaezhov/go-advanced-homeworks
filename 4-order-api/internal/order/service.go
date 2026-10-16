package order

import (
	"4-order-api/internal/product"
	"errors"
	"fmt"
)

type OrderService struct {
	OrderRepo   *OrderRepository
	ProductRepo *product.ProductRepository
}

func NewOrderService(orderRepo *OrderRepository, productRepo *product.ProductRepository) *OrderService {
	return &OrderService{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(userID uint, itemsReq []OrderItemRequest) (*Order, error) {
	if len(itemsReq) == 0 {
		return nil, errors.New("список товаров не может быть пустым")
	}

	var orderItems []OrderItem
	var totalPrice float64
	for _, itemReq := range itemsReq {
		prod, err := s.ProductRepo.Get(itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("товар с ID %d не найден: %w", itemReq.ProductID, err)
		}

		price := prod.Price
		totalPrice += price * float64(itemReq.Quantity)

		orderItems = append(orderItems, OrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     price,
		})
	}

	order := &Order{
		UserID: userID,
		Items:  orderItems,
		Price:  totalPrice,
	}

	return s.OrderRepo.Create(order)
}

func (s *OrderService) GetOrder(id uint) (*Order, error) {
	return s.OrderRepo.Get(id)
}

func (s *OrderService) ListMyOrders(userID uint, limit, offset int) (OrderListResponse, error) {
	return s.OrderRepo.List(OrderListPayload{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
}
