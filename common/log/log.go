package log

import (
	"context"

	"github.com/youthlin/glog/common/ctxs"
	"go.uber.org/zap"
)

// CtxAddKV add key-values to context
func CtxAddKV(ctx context.Context, kvs ...interface{}) context.Context {
	return context.WithValue(ctx, ctxs.CtxKeyLogKV, kvs)
}

func Debug(ctx context.Context, fmt string, args ...interface{}) {
	logger(ctx).Debugf(fmt, args...)
}
func Info(ctx context.Context, fmt string, args ...interface{}) {
	logger(ctx).Infof(fmt, args...)
}
func Warn(ctx context.Context, fmt string, args ...interface{}) {
	logger(ctx).Warnf(fmt, args...)
}
func Error(ctx context.Context, fmt string, args ...interface{}) {
	logger(ctx).Errorf(fmt, args...)
}

func logger(ctx context.Context) *zap.SugaredLogger {
	s := zap.L().WithOptions(zap.AddCallerSkip(1)).Sugar()
	if value, ok := ctx.Value(ctxs.CtxKeyLogKV).([]interface{}); ok {
		s = s.With(value...)
	}
	return s
}
