package order

import (
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/middleware"
	"4-order-api/pkg/req"
	"4-order-api/pkg/res"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type OrderHandlerDeps struct {
	OrderService *OrderService
	JWT          *jwt.JWT
}
type OrderHandler struct {
	OrderService *OrderService
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderService: deps.OrderService,
	}
	router.Handle("POST /order", middleware.IsAuthed(handler.Create(), deps.JWT))
	router.Handle("GET /order/{id}", middleware.IsAuthed(handler.Get(), deps.JWT))
	router.Handle("GET /my-orders", middleware.IsAuthed(handler.MyOrders(), deps.JWT))
}

func (handler *OrderHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			res.Json(w, "Ошибка получения пользователя", http.StatusBadRequest)
			return
		}

		id, err := req.GetIDFromPath(r)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		order, err := handler.OrderService.GetOrder(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "Заказ не найден", http.StatusNotFound)
				return
			}
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if order.UserID != userID {
			res.Json(w, "Заказ не найден", http.StatusNotFound)
			return
		}

		res.Json(w, order, http.StatusOK)
	}
}

func (handler *OrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			res.Json(w, "Ошибка получения пользователя", http.StatusBadRequest)
			return
		}

		body, err := req.HandleBody[OrderCreatePayload](&w, r)
		if err != nil {
			return
		}

		order, err := handler.OrderService.CreateOrder(userID, body.Items)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, order, http.StatusCreated)
	}
}

func (handler *OrderHandler) MyOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		limit := 10
		offset := 0

		if l := query.Get("limit"); l != "" {
			if val, err := strconv.Atoi(l); err == nil && val > 0 {
				limit = val
			}
		}
		if o := query.Get("offset"); o != "" {
			if val, err := strconv.Atoi(o); err == nil && val >= 0 {
				offset = val
			}
		}

		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			res.Json(w, "Ошибка получения пользователя", http.StatusBadRequest)
			return
		}

		ordersResponse, err := handler.OrderService.ListMyOrders(userID, limit, offset)

		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, ordersResponse, http.StatusOK)
	}
}
