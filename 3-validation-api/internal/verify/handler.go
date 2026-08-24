package verify

import (
	"3-validation-api/configs"
	"3-validation-api/internal/email"
	"3-validation-api/pkg/file"
	"3-validation-api/pkg/hash"
	"3-validation-api/pkg/req"
	"3-validation-api/pkg/res"
	"fmt"
	"net/http"
)

const FILE_NAME = "validate-db.json"

type VerifyHandlerDeps struct {
	*configs.Config
}
type VerifyHandler struct {
	*configs.Config
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := VerifyHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
}
func (handler *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[SendRequest](&w, r)
		if err != nil {
			return
		}

		fmt.Println(body)

		db, err := file.NewJSONFileStorage(FILE_NAME)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		vault := NewVault(db)

		key := hash.GenerateRandomHash()
		vault.Add(VaultEntry{Email: body.Email, Hash: key})

		to := body.Email
		link := fmt.Sprintf("%v/verify/%v", req.GetBaseURL(r), key)
		err = email.Send(handler.Config, to, link)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, fmt.Sprintf("Успешно отправлено на %v", to), http.StatusOK)
	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("hash")
		if key == "" {
			res.Json(w, "Ошибка валидации ссылки", http.StatusBadRequest)
			return
		}
		fmt.Println(key)

		db, err := file.NewJSONFileStorage(FILE_NAME)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		vault := NewVault(db)

		isDeleted := vault.Delete(key)

		if !isDeleted {
			res.Json(w, isDeleted, http.StatusBadRequest)
			return
		}

		res.Json(w, isDeleted, http.StatusOK)
	}
}
