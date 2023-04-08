package Utils

import (
	medianasms "github.com/medianasms/go-rest-sdk"
	"gopkg.in/gomail.v2"
)

var apiKey, _ = ReadFromEnvFile(".env", "API_KEY")
var googlePass, _ = ReadFromEnvFile(".env", "GOGLESECRET")

func SendOTP(phoneNumber string, otp string) error {
	sms := medianasms.New(apiKey)
	patternValue := map[string]string{
		"code": otp,
	}
	_, err := sms.SendPattern(
		"mxq9qfupb3xzcpz",
		"+985000125475",
		phoneNumber,
		patternValue)
	if err != nil {
		return err
	}
	return nil
}

func SendEmail(to string, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "mhmdrzsmip@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your OTP is here!")
	m.SetBody("text/plain", "کد تایید ورود به سامانه همینجا: "+otp)

	d := gomail.NewDialer("smtp.gmail.com", 465, "mhmdrzsmip@gmail.com", googlePass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
