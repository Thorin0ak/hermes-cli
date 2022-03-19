package token

import (
	"github.com/Thorin0ak/mercure-test/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewJWTMaker(t *testing.T) {
	m, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	sub := utils.RandomString(8)
	topic := "sse://foo.bar/tutu"
	duration := time.Minute
	pClaim := []string{topic}
	cMap := make(map[string][]string)
	cMap["publish"] = pClaim

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := m.CreateToken(sub, topic, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := m.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Subject, sub)
	require.WithinDuration(t, issuedAt, time.Unix(payload.IssuedAt, 0), time.Second)
	require.WithinDuration(t, expiredAt, time.Unix(payload.ExpiresAt, 0), time.Second)
	require.Equal(t, payload.Mercure, cMap)
}

func TestExpiredJWTToken(t *testing.T) {
	m, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	sub := utils.RandomString(8)
	topic := "sse://foo.bar/tutu"
	duration := -time.Minute

	token, err := m.CreateToken(sub, topic, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	p, err := m.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, p)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	p := NewPayload(utils.RandomString(8), "sse://foo.bar/tutu", time.Minute)
	require.NotEmpty(t, p)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, p)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	m, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	payload, err := m.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
