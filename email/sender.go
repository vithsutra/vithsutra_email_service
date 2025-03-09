package email

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/vithsutra/vithsutra_email_sending_service/internal/models"
)

func sendEmail(to string, subject string, htmlBody string) error {
	from := os.Getenv("ROOT_EMAIL")
	if from == "" {
		return errors.New("missing ROOT_EMAIL env variable")
	}

	password := os.Getenv("ROOT_EMAIL_PASSWORD")

	if password == "" {
		return errors.New("missing ROOT_EMAIL_PASSWORD env variable")
	}

	smtpHost := os.Getenv("SMTP_HOST")

	if smtpHost == "" {
		return errors.New("missing SMTP_HOST env variable")
	}

	smtpPort := os.Getenv("SMTP_PORT")

	if smtpPort == "" {
		return errors.New("missing SMTP_PORT env variable")
	}

	message := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	message += fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, htmlBody)

	err := smtp.SendMail(smtpHost+":"+smtpPort, smtp.PlainAuth("", from, password, smtpHost), from, []string{to}, []byte(message))

	return err
}

func SendOtpEmail(message *models.Email) error {
	template, err := parseTemplate("otp")

	if err != nil {
		return err
	}

	htmlTemplate, err := renderTemplate(template, message.Data)

	if err != nil {
		return err
	}

	if err := sendEmail(message.To, message.Subject, htmlTemplate); err != nil {
		return err
	}
	return nil
}

func WelcomeEmail(message *models.Email) error {
	template, err := parseTemplate("welcome")

	if err != nil {
		return err
	}

	htmlTemplate, err := renderTemplate(template, message.Data)

	if err != nil {
		return err
	}

	if err := sendEmail(message.To, message.Subject, htmlTemplate); err != nil {
		return err
	}
	return nil
}
