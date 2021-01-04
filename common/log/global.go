package log

import (
	"go.uber.org/zap"
)

var defaultLogger = NewLogger(1, zap.NewNop())

func SetLogger(zapL *zap.Logger) {
	defaultLogger = NewLogger(1, zapL)
}

func With(kvs ...interface{}) *Logger {
	return defaultLogger.with(-1, kvs...)
}

func Debug(fmt string, args ...interface{}) {
	defaultLogger.Debug(fmt, args...)
}
func Info(fmt string, args ...interface{}) {
	defaultLogger.Info(fmt, args...)
}
func Warn(fmt string, args ...interface{}) {
	defaultLogger.Warn(fmt, args...)
}
func Error(fmt string, args ...interface{}) {
	defaultLogger.Error(fmt, args...)
}
func Flush() {
	defaultLogger.Flush()
}
