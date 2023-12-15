package Utils

import (
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"testing"
)

func TestMedianaSMS_SendSMSWithPattern(t *testing.T) {
	sms := NewSMS(utils.ReadFromEnvFile(".env", "API_KEY"), map[string]string{
		"code": "1234",
	})
	err := sms.SendSMSWithPattern("09100570877", utils.ReadFromEnvFile(".env", "OTP_PATTERN_CODE"))
	if err != nil {
		t.Error(err)
	}
}
