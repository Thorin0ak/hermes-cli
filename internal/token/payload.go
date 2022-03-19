package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	Mercure map[string][]string `json:"mercure"`
	jwt.StandardClaims
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func (p *Token) Valid() error {
	if time.Now().After(time.Unix(p.ExpiresAt, 0)) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(sub string, topic string, duration time.Duration) *Token {
	pClaim := []string{topic}
	cMap := make(map[string][]string)
	cMap["publish"] = pClaim
	return &Token{
		cMap,
		jwt.StandardClaims{
			Issuer:    "test",
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   sub,
		},
	}
}
