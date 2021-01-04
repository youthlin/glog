package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/youthlin/glog/common"
	"github.com/youthlin/glog/common/ctxs"
	"github.com/youthlin/glog/common/log"
)

var background = context.Background()

func main() {
	ctx := context.WithValue(background, ctxs.CtxKeyStartedAt, time.Now())
	defer onExit(ctx)
	common.MustInit()
	config := common.Config()
	log.With("env", os.Environ()).
		Infoj("应用启动|App start|args=%v|app dir=%v|config=%#v", os.Args, common.AppDir(), config)
	log.Infoj("网络配置|Web config=%v", config.Web)
	gin.SetMode(config.Web.Mode)
	r := gin.Default()
	register(r)
	server := &http.Server{Addr: config.Web.Addr, Handler: r}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("监听异常将退出|listen error|err=%+v", err)
		}
	}()
	gracefulShutdown(server)
}

func onExit(ctx context.Context) {
	duration := time.Since(ctx.Value(ctxs.CtxKeyStartedAt).(time.Time))
	log.With("run", duration).Info("应用退出|App Exit")
	log.Flush()
}

func gracefulShutdown(server *http.Server) {
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	sig := <-signals
	log.Info("监听到信号，服务器将关闭|got a signal: [%v]", sig)
	ctx, cancel := context.WithTimeout(background, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("服务器关闭异常|server shutdown error|err=%+v", err)
		return
	}
	log.Info("服务器已关闭|server shutdown")
}
