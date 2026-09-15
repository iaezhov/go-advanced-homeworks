package product

import (
	"4-order-api/pkg/req"
	"4-order-api/pkg/res"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	ProductRepository *ProductRepository
}
type ProductHandler struct {
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
	}
	router.HandleFunc("POST /product", handler.Create())
	router.HandleFunc("PATCH /product/{id}", handler.Update())
	router.HandleFunc("DELETE /product/{id}", handler.Delete())
	router.HandleFunc("GET /product/{id}", handler.Get())
}

func (handler *ProductHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdFromPath(r)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		product, err := handler.ProductRepository.Get(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "Продукт не найден", http.StatusNotFound)
				return
			}
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, product, http.StatusOK)
	}
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ProductCreatePayload](&w, r)
		if err != nil {
			return
		}

		product, err := handler.ProductRepository.Create(&Product{
			Name:        body.Name,
			Description: body.Description,
			Price:       body.Price,
			Images:      body.Images,
		})
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, product, http.StatusCreated)
	}
}

func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ProductUpdatePayload](&w, r)
		if err != nil {
			return
		}

		id, err := getIdFromPath(r)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		product, err := handler.ProductRepository.Update(&Product{
			Model:       gorm.Model{ID: uint(id)},
			Name:        body.Name,
			Description: body.Description,
			Price:       body.Price,
			Images:      body.Images,
		})
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, product, http.StatusOK)
	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdFromPath(r)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.ProductRepository.Get(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "Продукт не найден", http.StatusNotFound)
				return
			}
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = handler.ProductRepository.Delete(id)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, nil, http.StatusOK)
	}
}

func getIdFromPath(r *http.Request) (uint, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return 0, errors.New("Отсутствует Id")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("Неверный формат id")
	}

	return uint(id), nil
}
