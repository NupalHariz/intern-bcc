package gomail

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendGoMail(otp string, email string) error {
	host := os.Getenv("GOMAIL_HOST")
	stringPort := os.Getenv("GOMAIL_PORT")
	port, _ := strconv.Atoi(stringPort)
	username := os.Getenv("GOMAIL_USERNAME")
	password := os.Getenv("GOMAIL_PASSWORD")
	body := fmt.Sprintf("This is your OTP code <b>%v</b> and <i>I'm Naufal</i>!", otp)

	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Test Send Email")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, username, password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
