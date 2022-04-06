package internal

import (
	"fmt"
	"github.com/Thorin0ak/mercure-test/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestNewOrchestratorOk(t *testing.T) {
	secret := utils.RandomString(32)
	hConf := GetConfig()
	envs := []MercureConfig{{
		Name:      "localhost",
		HubUrl:    "https://local/bar",
		JwtSecret: secret,
	}, {
		Name:      "dev",
		HubUrl:    "https://dev/bar",
		JwtSecret: secret,
	}}
	hConf.Mercure = &MercureEnvs{[]MercureConfig{}}
	hConf.Mercure.Envs = envs
	logger, _ := utils.GetObservedLogger(zap.InfoLevel)

	o, err := NewOrchestrator(hConf, logger)
	require.Empty(t, err)
	require.Equal(t, "https://local/bar", o.hubUrl)
}

func TestNewOrchestratorNoEnvMatch(t *testing.T) {
	// Not necessary, just playing with Zap
	var opts []zap.Option
	opts = append(opts, zap.OnFatal(zapcore.WriteThenPanic))
	observedZapCore, observedLogs := observer.New(zapcore.InfoLevel)
	logger := zap.New(observedZapCore, opts...).Sugar()

	secret := utils.RandomString(32)
	hConf := GetConfig()
	envs := []MercureConfig{{
		Name:      "test",
		HubUrl:    "https://test/bar",
		JwtSecret: secret,
	}, {
		Name:      "dev",
		HubUrl:    "https://dev/bar",
		JwtSecret: secret,
	}}
	hConf.Mercure = &MercureEnvs{[]MercureConfig{}}
	hConf.Mercure.Envs = envs

	o, err := NewOrchestrator(hConf, logger)
	require.Empty(t, o)
	require.EqualError(t, err, "no config found for active env localhost")
	require.Equal(t, 0, observedLogs.Len())
}

func TestNewOrchestratorInvalidSecret(t *testing.T) {
	secret := utils.RandomString(4)
	hConf := GetConfig()
	envs := []MercureConfig{{
		Name:      "localhost",
		HubUrl:    "https://local/bar",
		JwtSecret: secret,
	}}
	hConf.Mercure = &MercureEnvs{[]MercureConfig{}}
	hConf.Mercure.Envs = envs
	logger, _ := utils.GetObservedLogger(zap.InfoLevel)

	o, err := NewOrchestrator(hConf, logger)
	require.Empty(t, o)
	require.EqualError(t, err, fmt.Sprintf("invalid key size: must be at least %d chars", minSecretKeySize))
}

//type args struct {
//	config       *Config
//	logger       *zap.SugaredLogger
//	observedLogs *observer.ObservedLogs
//	errMsg       string
//}
//
//func newArgs(envs []MercureConfig, errMsg string) args {
//	hConf := GetConfig()
//	hConf.Mercure = &MercureEnvs{[]MercureConfig{}}
//	hConf.Mercure.Envs = envs
//
//	logger, observedLogs := utils.GetObservedLogger(zap.InfoLevel)
//
//	return args{
//		config:       hConf,
//		logger:       logger,
//		observedLogs: observedLogs,
//		errMsg:       errMsg,
//	}
//}
//
//func TestNewOrchestrator(t *testing.T) {
//	secret := utils.RandomString(34)
//
//	tests := []struct {
//		name    string
//		args    args
//		want    *Orchestrator
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name: "gets new orchestrator",
//			args: newArgs([]MercureConfig{{
//				Name:      "localhost",
//				HubUrl:    "https://local/bar",
//				JwtSecret: secret,
//			}, {
//				Name:      "dev",
//				HubUrl:    "https://dev/bar",
//				JwtSecret: secret,
//			}}, ""),
//			want:    GetDummy(nil, "https://local/bar", secret, nil),
//			wantErr: assert.NoError,
//		},
//		{
//			name: "fails to get new orchestrator with invalid secret",
//			args: newArgs([]MercureConfig{{
//				Name:      "localhost",
//				HubUrl:    "https://local/bar",
//				JwtSecret: utils.RandomString(4),
//			}}, fmt.Sprintf("invalid key size: must be at least %d chars", minSecretKeySize)),
//			want:    nil,
//			wantErr: assert.Error,
//		},
//		{
//			name: "fails to get new orchestrator with invalid active env",
//			args: newArgs([]MercureConfig{{
//				Name:      "foo",
//				HubUrl:    "https://local/bar",
//				JwtSecret: secret,
//			}}, "no config found for active env localhost"),
//			want:    nil,
//			wantErr: assert.Error,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := NewOrchestrator(tt.args.config, tt.args.logger)
//			if !tt.wantErr(t, err, tt.args.errMsg) {
//				return
//			}
//			assert.Equalf(t, tt.want, got, "NewOrchestrator(%v, %v)", tt.args.config, tt.args.logger)
//		})
//	}
//}
