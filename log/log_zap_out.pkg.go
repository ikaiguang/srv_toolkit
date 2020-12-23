package tklog

import (
	"context"
	"go.uber.org/zap"
)

// Info logs a message at the info log level.
func Info(msg string) {
	logger.Info(msg)
}

// Error logs a message at the error log level.
func Error(msg string) {
	logger.Error(msg)
}

// Warn logs a message at the warning log level.
func Warn(msg string) {
	logger.Warn(msg)
}

// Infov logs a message at the info log level.
func Infov(args ...zap.Field) {
	logger.Info("", args...)
}

// Errorv logs a message at the error log level.
func Errorv(args ...zap.Field) {
	logger.Error("", args...)
}

// Warnv logs a message at the warning log level.
func Warnv(args ...zap.Field) {
	logger.Warn("", args...)
}

// Infoc logs a message at the info log level.
func Infoc(ctx context.Context, msg string, args ...zap.Field) {
	args = append(args, AddExtraField(ctx)...)
	logger.Info(msg, args...)
}

// Errorc logs a message at the error log level.
func Errorc(ctx context.Context, msg string, args ...zap.Field) {
	args = append(args, AddExtraField(ctx)...)
	logger.Error(msg, args...)
}

// Warnc logs a message at the warning log level.
func Warnc(ctx context.Context, msg string, args ...zap.Field) {
	args = append(args, AddExtraField(ctx)...)
	logger.Warn(msg, args...)
}
