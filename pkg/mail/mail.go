package mail

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v6"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Config struct {
	ClientID     string `env:"GMAIL_CLIENT_ID"`
	ClientSecret string `env:"GMAIL_CLIENT_SECRET"`
	AccessToken  string `env:"GMAIL_ACCESS_TOKEN"`
	RefreshToken string `env:"GMAIL_REFRESH_TOKEN"`
}

type Mailer struct {
	Config       Config
	GMailService *gmail.Service
}

func NewMailer() (Mailer, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Mailer{}, err
	}

	log.Printf("Config: %+v\n", cfg)

	mailer := Mailer{
		Config:       cfg,
		GMailService: oAuthGmailService(cfg),
	}

	return mailer, nil
}

func oAuthGmailService(c Config) *gmail.Service {
	config := oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  c.AccessToken,
		RefreshToken: c.RefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)
	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	if srv != nil {
		fmt.Println("Email service is initialized.")
	}

	return srv
}

func (m *Mailer) SendEmailSMTP(to string, subject string, data interface{}, template string) error {
	emailBody, err := parseTemplate(template, data)
	if err != nil {
		return errors.New("unable to parse email template")
	}

	var message gmail.Message
	emailTo := "To: " + to + "\r\n"
	emailSubject := "Subject: " + subject + "\n"
	mime := "MIME-version: 1.0\nContent-Type: text/plain; charset:\"UTF-8\";\n\n"
	msg := []byte(emailTo + emailSubject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = m.GMailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
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
