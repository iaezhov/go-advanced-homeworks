package order

type OrderItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

type OrderCreatePayload struct {
	Items []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type MyOrdersPayload struct {
	Limit  int
	Offset int
}

type OrderListPayload struct {
	UserID uint
	Limit  int
	Offset int
}

type OrderListResponse struct {
	Items []Order
	Total int
}
