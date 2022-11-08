package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	// Initialize a new mail dialer which will given the SMTP server settings
	// as well as a timeout to be used whenever an email is sent.
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, templateFileName string, data any) error {
	// Initialize the template to be used for emails from the file system.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFileName)
	if err != nil {
		return err
	}

	// Write the subject template a byte buffer with any dynamic data.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Write the plainBody template a byte buffer with any dynamic data.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// Write the htmlBody template a byte buffer with any dynamic data.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Create a message to be sent and set the dynamic template data accordingly.
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.SetBody("text/html", htmlBody.String())

	// Try sending the email message 3 times. If it fails, return the error.
	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)
		// Err is Nil so the email has succeeded; return.
		if nil == err {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return err
}
