package product

import (
	"4-order-api/pkg/db"
	"4-order-api/pkg/res"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	Db *db.Db
}
type ProductHandler struct {
	db *db.Db
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		db: deps.Db,
	}
	router.HandleFunc("GET /product/{id}", handler.GetProduct())
}

func (handler *ProductHandler) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			res.Json(w, "Ошибка валидации ссылки", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil || id == 0 {
			res.Json(w, "Неверный формат ID", http.StatusBadRequest)
			return
		}

		productService := NewProductService(handler.db)

		product, err := productService.GetByID(uint(id))

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "Продукт не найден", http.StatusNotFound)
				return
			}
			res.Json(w, "Ошибка базы данных", http.StatusInternalServerError)
			return
		}

		res.Json(w, product, http.StatusOK)
	}
}
