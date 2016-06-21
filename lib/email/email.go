package email

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
)

var (
	e SMTPInfo
)

// SMTPInfo is the details for the SMTP server
type SMTPInfo struct {
	Username string
	Password string
	Hostname string
	Port     int
	From     string
}

// SetConfig adds the settings for the SMTP server
func SetConfig(c SMTPInfo) {
	e = c
}

// Config returns the configuration
func Config() SMTPInfo {
	return e
}

// Send mails an email
func Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Hostname)

	header := make(map[string]string)
	header["From"] = e.From
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = `text/plain; charset="utf-8"`
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.Hostname, e.Port),
		auth,
		e.From,
		[]string{to},
		[]byte(message),
	)

	return err
}
