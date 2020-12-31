package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/youthlin/glog/common"
	"github.com/youthlin/glog/common/ctxs"
	"github.com/youthlin/glog/common/log"
	"go.uber.org/zap"
)

func main() {
	ctx := context.WithValue(context.Background(), ctxs.CtxKeyStartedAt, time.Now())
	defer onExit(ctx)
	common.MustInit()
	config := common.Config()
	log.Info(log.CtxAddKV(ctx, "env", os.Environ()), "App start|config=%#v", config)

	r := gin.Default()
	gin.SetMode(config.Web.Mode)
	register(r)
	log.Info(ctx, "Run: %+v", r.Run(config.Web.Addr))
}

func onExit(ctx context.Context) {
	log.Info(log.CtxAddKV(ctx, "run", time.Since(ctx.Value(ctxs.CtxKeyStartedAt).(time.Time))), "App Exit")
	_ = zap.L().Sync()
}
