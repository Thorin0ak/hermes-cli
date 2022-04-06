package internal

import (
	"fmt"
	"github.com/Thorin0ak/hermes-cli/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestNewJWTMaker(t *testing.T) {
	logger, _ := utils.GetObservedLogger(zap.InfoLevel)
	m, err := NewJWTMaker(utils.RandomString(32), logger)
	require.NoError(t, err)

	sub := utils.RandomString(8)
	topic := "sse://foo.bar/tutu"
	duration := time.Minute
	pClaim := []string{topic}
	cMap := make(map[string][]string)
	cMap["publish"] = pClaim

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := m.CreateToken(sub, topic, duration, true)
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
	logger, _ := utils.GetObservedLogger(zap.InfoLevel)
	m, err := NewJWTMaker(utils.RandomString(32), logger)
	require.NoError(t, err)

	sub := utils.RandomString(8)
	topic := "sse://foo.bar/tutu"
	duration := -time.Minute

	token, err := m.CreateToken(sub, topic, duration, true)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	p, err := m.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, p)
}

func TestSubscribeJTWTToken(t *testing.T) {
	logger, _ := utils.GetObservedLogger(zap.InfoLevel)
	m, err := NewJWTMaker(utils.RandomString(32), logger)
	require.NoError(t, err)

	sub := utils.RandomString(8)
	topic := "sse://foo.bar/tutu"
	duration := time.Minute
	isPublisher := false

	token, err := m.CreateToken(sub, topic, duration, isPublisher)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := m.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	customClaims := map[string][]string{
		"subscribe": {topic},
	}
	require.Equal(t, payload.Mercure, customClaims)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	p := NewPayload(utils.RandomString(8), "sse://foo.bar/tutu", time.Minute, true)
	require.NotEmpty(t, p)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, p)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	logger, _ := utils.GetObservedLogger(zap.InfoLevel)
	m, err := NewJWTMaker(utils.RandomString(32), logger)
	require.NoError(t, err)

	payload, err := m.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestInvalidSecretKey(t *testing.T) {
	logger, _ := utils.GetObservedLogger(zap.InfoLevel)
	m, err := NewJWTMaker(utils.RandomString(2), logger)
	errMsg := fmt.Sprintf("invalid key size: must be at least %d chars", minSecretKeySize)
	require.EqualErrorf(t, err, errMsg, "")
	require.Nil(t, m)
}
