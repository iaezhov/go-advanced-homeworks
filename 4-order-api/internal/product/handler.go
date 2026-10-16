package product

import (
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/middleware"
	"4-order-api/pkg/req"
	"4-order-api/pkg/res"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	ProductRepository *ProductRepository
	JWT               *jwt.JWT
}
type ProductHandler struct {
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
	}
	router.Handle("POST /product", middleware.IsAuthed(handler.Create(), deps.JWT))
	router.Handle("PATCH /product/{id}", middleware.IsAuthed(handler.Update(), deps.JWT))
	router.Handle("DELETE /product/{id}", middleware.IsAuthed(handler.Delete(), deps.JWT))
	router.Handle("GET /product/{id}", middleware.IsAuthed(handler.Get(), deps.JWT))
}

func (handler *ProductHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := req.GetIDFromPath(r)
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

		id, err := req.GetIDFromPath(r)
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
		id, err := req.GetIDFromPath(r)
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
