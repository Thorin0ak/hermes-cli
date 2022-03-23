package internal

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type Maker interface {
	CreateToken(sub string, topic string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Token, error)
}

type JWTMaker struct {
	secretKey string
}

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
	//cMap["publish"] = pClaim
	cMap["subscribe"] = pClaim
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
