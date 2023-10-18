package Utils

import medianaSMS "github.com/medianasms/go-rest-sdk"

type SMSService interface {
	SendSMSWithPattern(phoneNumber string, patternCode string) error
}

type MedianaSMS struct {
	ApiKey  string
	Pattern map[string]string
}

func NewSMS(apiKey string, pattern map[string]string) *MedianaSMS {
	return &MedianaSMS{
		ApiKey:  apiKey,
		Pattern: pattern,
	}
}

func (sms *MedianaSMS) SendSMSWithPattern(phoneNumber string, patternCode string) error {
	smsService := medianaSMS.New(sms.ApiKey)
	_, err := smsService.SendPattern(
		patternCode, "+985000125475", phoneNumber, sms.Pattern)
	if err != nil {
		return err
	}
	return nil
}
