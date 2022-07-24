package main

import (
	"flag"
	"log"

	"github.com/ghiblin/go-gmail/pkg/mail"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	to := flag.String("to", "", "address to send email to")
	flag.Parse()

	mailer, err := mail.NewMailer()
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		ReceiverName string
		SenderName   string
	}{
		ReceiverName: "Alessandro",
		SenderName:   "Golang",
	}

	err = mailer.SendEmailSMTP([]string{*to}, "Test email", data, "sample_template.txt")
	if err != nil {
		log.Fatal(err)
	}
}
