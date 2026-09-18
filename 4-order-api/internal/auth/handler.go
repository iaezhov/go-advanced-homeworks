package auth

import (
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/req"
	"4-order-api/pkg/res"
	"errors"
	"net/http"
)

type AuthHandlerDeps struct {
	AuthService *AuthService
	JWT         *jwt.JWT
}
type AuthHandler struct {
	AuthService *AuthService
	JWT         *jwt.JWT
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		AuthService: deps.AuthService,
		JWT:         deps.JWT,
	}
	router.HandleFunc("POST /login", handler.Login())
	router.HandleFunc("POST /verify", handler.Verify())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginPayload](&w, r)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		sessionId, err := handler.AuthService.Login(body.Phone)
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := LoginResponse{
			SessionId: sessionId,
		}

		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[VerifyPayload](&w, r)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		phone, err := handler.AuthService.Verify(body.SessionId, body.Code)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrWrongCredentials) {
				statusCode = http.StatusUnauthorized
			}
			res.Json(w, err.Error(), statusCode)
			return
		}

		token, err := handler.JWT.Create(phone)
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := VerifyResponse{
			Token: token,
		}

		res.Json(w, data, http.StatusOK)
	}
}
