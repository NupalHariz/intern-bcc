package gomail

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type IGoMail interface {
	SendGoMail(subject string, htmlBody string, toEmail string) error
}

type Gomail struct {
	host     string
	port     int
	username string
	password string
}

func GoMailInit() IGoMail {
	host := os.Getenv("GOMAIL_HOST")
	stringPort := os.Getenv("GOMAIL_PORT")
	port, _ := strconv.Atoi(stringPort)
	username := os.Getenv("GOMAIL_USERNAME")
	password := os.Getenv("GOMAIL_PASSWORD")

	return &Gomail{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (g *Gomail) SendGoMail(subject string, htmlBody string, toEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", g.username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(g.host, g.port, g.username, g.password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
