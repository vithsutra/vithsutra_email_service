//boiler plate code

package email

import (
	"bytes"
	"email-service/config"
	"email-service/internal/models"
	"fmt"
	"html/template"
	"net/smtp"
)

func SendEmail(data models.EmailRequest) error {
	// Load template based on Type
	tmpl, err := template.ParseFiles("templates/" + data.Type + ".html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return err
	}

	// SMTP Authentication
	auth := smtp.PlainAuth("", config.Config.SMTP.Username, config.Config.SMTP.Password, config.Config.SMTP.Host)

	// Email Message
	msg := []byte("Subject: " + data.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		body.String())

	// Sending Email
	addr := fmt.Sprintf("%s:%d", config.Config.SMTP.Host, config.Config.SMTP.Port)
	err = smtp.SendMail(addr, auth, config.Config.SMTP.From, []string{data.To}, msg)
	if err != nil {
		return err
	}

	fmt.Println("Email sent successfully to:", data.To)
	return nil
}
