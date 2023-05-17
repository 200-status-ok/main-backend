package Token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strconv"
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

func (j *JWTMaker) MakeToken(username string, userID uint64, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, userID, role, duration)
	if err != nil {
		return "", nil, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": payload.Username,
		"id":       payload.ID,
		"userId":   payload.UserID,
		"role":     payload.Role,
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
	id, _ := uuid.Parse(claims["id"].(string))
	issuedAtStr := claims["issuedAt"].(string)
	issuedAt, err := time.Parse(time.RFC3339, issuedAtStr)
	expiredStr := claims["expired"].(string)
	expired, err := time.Parse(time.RFC3339, expiredStr)
	userID, _ := strconv.ParseUint(fmt.Sprintf("%v", claims["userId"].(float64)), 10, 64)

	payload := &Payload{
		Username: claims["username"].(string),
		ID:       id,
		UserID:   userID,
		IssuedAt: issuedAt,
		Expired:  expired,
	}
	if err := payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil

}
