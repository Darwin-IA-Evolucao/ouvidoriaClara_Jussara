package config

import (
	"os"

	"gopkg.in/gomail.v2"
)

func EnviarEmail(destinatario, assunto, contentType, mensagem string) error {
	SENHA_EMAIL := os.Getenv("EMAIL_PASS")
	EMAIL_REMETENTE := os.Getenv("EMAIL_FROM")
	m := gomail.NewMessage()
	m.SetHeader("From", EMAIL_REMETENTE)
	m.SetHeader("To", destinatario)
	m.SetHeader("Subject", assunto)
	m.SetBody(contentType, mensagem)

	d := gomail.NewDialer("smtp.gmail.com", 587, EMAIL_REMETENTE, SENHA_EMAIL)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
