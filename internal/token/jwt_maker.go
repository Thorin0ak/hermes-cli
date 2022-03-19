package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32

func (m *JWTMaker) CreateToken(sub string, topic string, duration time.Duration) (string, error) {
	p := NewPayload(sub, topic, duration)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, p)
	return jwtToken.SignedString([]byte(m.secretKey))
}

func (m *JWTMaker) VerifyToken(token string) (*Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Printf("error determining signature algorithm\n")
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Token{}, keyFunc)
	if err != nil {
		fmt.Printf("error parsing token\n")
		validationError, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(validationError.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Token)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d chars", minSecretKeySize)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}
