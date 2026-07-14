package verify

import (
	"3-validation-api/3-validation-api/configs"
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/jordan-wright/email"
)

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
		fmt.Print("Send")
		e := email.NewEmail()
		e.From = "Jordan Wright <test@gmail.com>"
		e.To = []string{"test@example.com"}
		e.Subject = "Awesome Subject"
		e.Text = []byte("Text Body is, of course, supported!")
		e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
		e.Send("smtp.gmail.com:587", smtp.PlainAuth("", handler.Config.Email, handler.Config.Password, handler.Config.Address))

	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("Verify")
	}
}
