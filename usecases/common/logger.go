package common

import (
  "go.uber.org/zap"
)

type Logger interface {
	Info(message string, tags ...zap.Field)
	Panic(message string, err error, tags ...zap.Field)
	Fatal(message string, err error, tags ...zap.Field)
	Warn(message string, tags ...zap.Field)
	Error(message string, err error, tags ...zap.Field)
	Sync()
	Debug(message string, tags ...zap.Field)
}