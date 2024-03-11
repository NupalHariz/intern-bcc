package gomail

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type IGoMail interface{
	SendGoMail(otp string, email string) error
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
		host: host,
		port: port,
		username: username,
		password: password,
	}
}

func (g *Gomail) SendGoMail(otp string, email string) error {
	body := fmt.Sprintf("This is your OTP code <b>%v</b> and <i>I'm Naufal</i>!", otp)

	m := gomail.NewMessage()
	m.SetHeader("From", g.username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Test Send Email")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(g.host, g.port, g.username, g.password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
