package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"os"
	"text/template"

	"github.com/fakovacic/aig/internal/log"
	"gopkg.in/gomail.v2"
)

// New return new mailer
func New() *Mailer {
	// TODO: move to settings from envs
	d := gomail.NewDialer(
		os.Getenv("EMAIL_HOST"),
		587,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASS"),
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	dial, err := d.Dial()
	if err != nil {
		log.Errorf("dial email client error:%v", err)
	}
	dial.Close()
	return &Mailer{
		d: d,
	}
}

// Mailer struct
type Mailer struct {
	d *gomail.Dialer
}

// Send sends email
func (d Mailer) Send(to, subject, tplName string, data interface{}) error {
	tpl, err := d.ParseTemplate(tplName, data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", os.Getenv("EMAIL_USER"))
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", tpl)
	err = d.d.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

// ParseTemplate parse email template
func (d Mailer) ParseTemplate(tplName string, data interface{}) (string, error) {
	t := template.New("email.html")
	var err error
	t, err = t.ParseFiles(
		"templates/email/email.html",
		fmt.Sprintf("templates/email/%s.html", tplName))
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
