// Package email provides email sending via SMTP.
package email

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
)

// Info holds the details for the SMTP server.
type Info struct {
	Username string
	Password string
	Hostname string
	Port     int
	From     string
}

// header generates a plaintext header for the email.
func header(c Info, to, subject, body string) map[string]string {
	// Create the header
	header := map[string]string{
		"From":                      c.From,
		"To":                        to,
		"Subject":                   subject,
		"MIME-Version":              "1.0",
		"Content-Type":              `text/html; charset="utf-8"`,
		"Content-Transfer-Encoding": "base64",
	}

	return header
}

// Send an email.
func (c Info) Send(to, subject, body string) error {
	// Authentication for SMTP
	auth := smtp.PlainAuth("", c.Username, c.Password, c.Hostname)

	// Create header
	header := header(c, to, subject, body)

	// Set the message
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", c.Hostname, c.Port),
		auth,
		c.From,
		[]string{to},
		[]byte(message),
	)

	return err
}
