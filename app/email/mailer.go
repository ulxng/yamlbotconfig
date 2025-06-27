package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type Mailer struct {
	conf SmtpConfig
}

func NewMailer(conf SmtpConfig) *Mailer {
	return &Mailer{conf: conf}
}

func (m Mailer) Send(message, subject string) error {
	e := email.NewEmail()
	e.From = m.conf.DefaultEmailFrom
	e.To = []string{m.conf.DefaultEmailTo}
	e.Subject = subject
	e.Text = []byte(message)

	addr := fmt.Sprintf("%s:%d", m.conf.Host, m.conf.Port)
	auth := smtp.PlainAuth("", m.conf.Username, m.conf.Password, m.conf.Host)
	tlsCfg := &tls.Config{ServerName: m.conf.Host}

	//для 465 хардкод, но так пока надо
	if err := e.SendWithTLS(addr, auth, tlsCfg); err != nil {
		return fmt.Errorf("email send (SMTPS 465): %w", err)
	}
	return nil
}
