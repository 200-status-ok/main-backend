package Utils

import (
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

var accountSid, _ = ReadFromEnvFile(".env", "TWILIO_ACCOUNT_SID")
var authToken, _ = ReadFromEnvFile(".env", "TWILIO_AUTH_TOKEN")
var fromNumber, _ = ReadFromEnvFile(".env", "TWILIO_FROM_NUMBER")

var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: accountSid,
	Password: authToken,
})

func SendOTP(phoneNumber string, otp string) error {
	params := &api.CreateMessageParams{}
	params.SetFrom(fromNumber)
	params.SetTo(phoneNumber)
	params.SetBody("Your verification code is: " + otp)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}
	return nil
}
