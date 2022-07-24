package mail

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host     string `env:"EMAIL_HOST"`
	From     string `env:"EMAIL_FROM"`
	Password string `env:"EMAIL_PASSWORD"`
	Port     int    `env:"EMAIL_PORT" envDefault:"587"`
}

type Mailer struct {
	Config Config
}

func NewMailer() (Mailer, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Mailer{}, err
	}

	log.Printf("Config: %+v\n", cfg)

	mailer := Mailer{
		Config: cfg,
	}

	return mailer, nil
}

func (m *Mailer) SendEmailSMTP(to []string, subject string, data interface{}, template string) error {
	emailAuth := smtp.PlainAuth("", m.Config.From, m.Config.Password, m.Config.Host)
	emailBody, err := parseTemplate(template, data)
	if err != nil {
		return errors.New("unable to parse email template")
	}

	mime := "Subject: " + subject + "\nMIME-version: 1.0\nContent-Type: text/plain; charset:\"UTF-8\";\n\n"
	msg := []byte(mime + "\n" + emailBody)
	addr := fmt.Sprintf("%s:%d", m.Config.Host, m.Config.Port)

	if err := smtp.SendMail(addr, emailAuth, m.Config.From, to, msg); err != nil {
		return err
	}
	return nil
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("templates/%s", templateFileName))
	if err != nil {
		return "", errors.New("invalid template name")
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}
