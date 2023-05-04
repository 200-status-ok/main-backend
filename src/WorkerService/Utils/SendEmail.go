package Utils

import "gopkg.in/gomail.v2"

type EmailService interface {
	SendEmailWithGoogle() error
}

type GoogleEmail struct {
	Service        *gomail.Message
	GooglePassword string
	From           string
	To             string
	Subject        string
	Body           string
}

func NewEmail(from string, to string, subject string, body string, googlePass string) *GoogleEmail {
	return &GoogleEmail{
		Service:        gomail.NewMessage(),
		From:           from,
		To:             to,
		Subject:        subject,
		Body:           body,
		GooglePassword: googlePass,
	}
}

func (email *GoogleEmail) SendEmailWithGoogle() error {
	email.Service.SetHeader("From", email.From)
	email.Service.SetHeader("To", email.To)
	email.Service.SetHeader("Subject", email.Subject)
	email.Service.SetBody("text/plain", email.Body)

	d := gomail.NewDialer("smtp.gmail.com", 465, email.From, email.GooglePassword)

	if err := d.DialAndSend(email.Service); err != nil {
		return err
	}
	return nil
}
