package ctxs

type CtxKey string

const (
	CtxKeyLogKV     CtxKey = "logKV"     // 需要输出到日志中的额外的键值对 value 是 ...interface{}
	CtxKeyStartedAt CtxKey = "startedAt" // 启动时间
)
