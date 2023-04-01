package Token

type Maker interface {
	MakeToken(username string, duration int64) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
