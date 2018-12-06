package mail

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"net/smtp"
)

type UpdateMail struct {
	From   string
	To     []string
	Data   string
	Update *Update
	Auth   smtp.Auth
}

func New(to []string, update *Update, auth smtp.Auth) *UpdateMail {
	return &UpdateMail{
		To:     to,
		Update: update,
		From:   fmt.Sprintf("%s", viper.GetString("smtp_from")),
		Auth:   auth,
	}
}

func (u *UpdateMail) Send() error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := fmt.Sprintf("Subject: [SUBSTITUTES] Class %s update for %s\n", u.Update.Class, u.Update.Date)
	final := []byte(subject + mime + "\n" + u.Data)
	return smtp.SendMail(fmt.Sprintf("%s:%v", viper.GetString("smtp_host"), viper.GetString("smtp_port")),
		u.Auth,
		u.From, u.To, final)
}

func (u *UpdateMail) Parse(file string) error {
	templ, err := template.ParseFiles(file)
	if err != nil {
		return err
	}

	b := new(bytes.Buffer)
	if err := templ.Execute(b, u.Update); err != nil {
		return err
	}
	u.Data = b.String()
	return nil
}
