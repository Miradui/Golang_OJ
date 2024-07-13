package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Get <3099349352@qq.com>"
	e.To = []string{"18850467248@163.com"}
	e.Subject = "Test"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")

	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "3099349352@qq.com",
		"xggjrdhiurdwdedg", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		t.Fatal(err)
	}
}
