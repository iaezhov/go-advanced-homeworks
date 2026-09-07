package email

import (
	"3-validation-api/configs"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func Send(config *configs.Config, to string, link string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%v <%v>", "Jordan Wright", config.Email)
	e.To = []string{to}
	e.Subject = "Подтверждение Email"
	e.Text = []byte("Подтверждение Email")
	e.HTML = fmt.Appendf(nil, `<a href="%v">Подтвердить Email</a>`, link)
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", config.Email, config.Password, config.Address))
	return err
}
