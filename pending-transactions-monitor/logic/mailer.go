package pendingTransactionsMonitor

import (
	"fmt"
	"net/smtp"
	"net/url"
)

// Letter represents an email latter
type Letter struct {
	To      string
	Subject string
	Body    string
}

// Mailer represents a mail sender
type Mailer struct {
	host      string
	username  string
	password  string
	toAddress string
}

// NewMailer creates a new instance of the Mail
func NewMailer(host string, username string, password string, toAddress string) Mailer {
	return Mailer{
		host:      host,
		username:  username,
		password:  password,
		toAddress: toAddress,
	}
}

// SendMail sends a letter
func (m Mailer) SendMail(l *Letter) error {
	host, err := url.Parse("//" + m.host)
	if err != nil {
		fmt.Println("Mailer.SendMail > Error (url.Parse): host, username, toAddress ", m.host, m.username, m.toAddress, "\n", err)
		return err
	}

	auth := smtp.PlainAuth("", m.username, m.password, host.Hostname())

	to := []string{l.To}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		mime+"\r\n"+
		"\r\n"+
		"%s\r\n", l.To, m.username, l.Subject, l.Body)
	msg := []byte(body)
	err = smtp.SendMail(m.host, auth, m.username, to, msg)

	if err != nil {
		fmt.Println("Mailer.SendMail > Error (smtp.SendMail): host, username, toAddress ", m.host, m.username, m.toAddress, "\n", err)
	}

	return err
}
