package Token

import "time"

type Maker interface {
	MakeToken(username string, userID uint64, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
