package token

import "time"

type Maker interface {
	CreateToken(sub string, topic string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Token, error)
}
