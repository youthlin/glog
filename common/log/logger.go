package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	base   *zap.Logger
	logger *zap.SugaredLogger
}

// NewLogger create a logger, skip: 相对于本结构体跳过几层调用
func NewLogger(skip int, zapL *zap.Logger) *Logger {
	zapL = zapL.WithOptions(zap.AddCallerSkip(skip + 1))
	return &Logger{
		base:   zapL,
		logger: zapL.Sugar(),
	}
}

func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.logger.Debugf(fmt, args...)
}
func (l *Logger) Info(fmt string, args ...interface{}) {
	l.logger.Infof(fmt, args...)
}
func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.logger.Warnf(fmt, args...)
}
func (l *Logger) Error(fmt string, args ...interface{}) {
	l.logger.Errorf(fmt, args...)
}

func (l *Logger) With(kvs ...interface{}) *Logger {
	s := l.logger.With(kvs...)
	return &Logger{
		base:   s.Desugar(),
		logger: s,
	}
}

func (l *Logger) Debugj(fmt string, args ...interface{}) {
	if l.enable(zapcore.DebugLevel) {
		l.logger.Debugf(fmt, toJson(args)...)
	}
}
func (l *Logger) Infoj(fmt string, args ...interface{}) {
	if l.enable(zapcore.InfoLevel) {
		l.logger.Infof(fmt, toJson(args)...)
	}
}
func (l *Logger) Warnj(fmt string, args ...interface{}) {
	if l.enable(zapcore.WarnLevel) {
		l.logger.Warnf(fmt, toJson(args)...)
	}
}
func (l *Logger) Errorj(fmt string, args ...interface{}) {
	if l.enable(zapcore.ErrorLevel) {
		l.logger.Errorf(fmt, toJson(args)...)
	}
}
func (l *Logger) Flush() {
	_ = l.base.Sync()
}
func (l *Logger) with(skip int, kvs ...interface{}) *Logger {
	base := l.base.WithOptions(zap.AddCallerSkip(skip))
	sug := base.Sugar().With(kvs...)
	return &Logger{
		base:   base,
		logger: sug,
	}
}

func (l *Logger) enable(level zapcore.Level) bool {
	return l.base.Core().Enabled(level)
}
