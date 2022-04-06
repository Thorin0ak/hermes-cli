package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func GetObservedLogger(level zapcore.Level) (*zap.SugaredLogger, *observer.ObservedLogs) {
	observedZapCore, observedLogs := observer.New(level)
	return zap.New(observedZapCore).Sugar(), observedLogs
}
