package Utils

import (
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"testing"
)

func TestGoogleEmail_SendEmailWithGoogle(t *testing.T) {
	email := NewEmail("mhmdrzsmip@gmail.com", "alifakhary622@gmail.com",
		"test", "test", utils.ReadFromEnvFile(".env", "GOOGLE_SECRET"))
	err := email.SendEmailWithGoogle()
	if err != nil {
		t.Error(err)
	}
}
