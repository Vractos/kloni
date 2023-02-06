package metrics

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

func NewLogger(logLevel string) *Logger {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(getLogLevel(logLevel)),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:      "level",
			TimeKey:       "time",
			MessageKey:    "message",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
	}

	log, err := logConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Panic(err.Error())
	}
	return &Logger{log: log}
}

func (l *Logger) Info(message string, tags ...zap.Field) {
	l.log.Info(message, tags...)
}

func (l *Logger) Panic(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.log.Panic(message, tags...)
}

func (l *Logger) Fatal(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.log.Fatal(message, tags...)
}

func (l *Logger) Warn(message string, tags ...zap.Field) {
	l.log.Warn(message, tags...)
}

func (l *Logger) Error(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.log.Error(message, tags...)
}

func (l *Logger) Sync() {
	l.log.Sync()
}

func (l *Logger) Debug(message string, tags ...zap.Field) {
	l.log.Debug(message, tags...)
}

func getLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}
