package Utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func EmailRandomGenerator() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	emailLength := 10
	email := ""

	for i := 0; i < emailLength; i++ {
		randomIndex := rand.Intn(len(chars))
		randomChar := chars[randomIndex]

		email += string(randomChar)
	}
	email += "@gmail.com"

	return email
}
