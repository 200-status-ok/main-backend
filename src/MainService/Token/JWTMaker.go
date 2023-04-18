package Token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, errors.New("secret key is too short")
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

func (j *JWTMaker) MakeToken(username string, duration int64) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": payload.Username,
		"id":       payload.ID,
		"issuedAt": payload.IssuedAt,
		"expired":  payload.Expired,
	})
	token, err := jwtToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	payload := &Payload{
		Username: claims["username"].(string),
		ID:       claims["id"].(uuid.UUID),
		IssuedAt: claims["issuedAt"].(time.Time),
		Expired:  claims["expired"].(time.Time),
	}
	if err := payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil

}
