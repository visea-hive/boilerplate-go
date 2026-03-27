package mail

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// Config holds the email composition details.
type Config struct {
	To         []string
	Cc         []string
	Bcc        []string
	Subject    string
	Body       string
	Attachment *bytes.Buffer
}

// Send delivers an email using SMTP settings from environment variables.
// Required env vars: MAIL_HOST, MAIL_PORT, MAIL_FROM_ADDRESS, MAIL_PASSWORD.
func Send(to, subject, body string) error {
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPortStr := os.Getenv("MAIL_PORT")
	from := os.Getenv("MAIL_FROM_ADDRESS")
	password := os.Getenv("MAIL_PASSWORD")

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("failed to convert MAIL_PORT to integer: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
