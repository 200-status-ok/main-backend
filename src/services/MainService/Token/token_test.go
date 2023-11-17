package Token

import (
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	secretKey := "Qwertyuiopqwertyuiopasdfghjkl123456789"
	maker, err := NewJWTMaker(secretKey)
	if err != nil {
		t.Errorf("failed to create JWTMaker: %s", err.Error())
	}
	username := "username"
	userID := uint64(1)
	role := "admin"
	duration := time.Duration(1) * time.Minute

	token, payload, err := maker.MakeToken(username, userID, role, duration)
	if err != nil {
		t.Errorf("failed to make token: %s", err.Error())
	}
	verifiedPayload, err := maker.VerifyToken(token)
	if err != nil {
		t.Errorf("failed to verify token: %s", err.Error())
	}
	if verifiedPayload.Username != payload.Username {
		t.Errorf("username is not the same")
	}
}
