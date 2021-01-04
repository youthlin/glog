package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/youthlin/glog/common/log"
	"github.com/youthlin/glog/common/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 在Go语言项目中使用Zap日志库 https://www.liwenzhou.com/posts/Go/zap/
func initLog() {
	configs := Config().Log
	var core []zapcore.Core
	for i := range configs {
		config := &configs[i]
		if config.Enable {
			// zapcore.Encoder + zapcore.WriteSyncer => zapcore.Core
			core = append(core, zapcore.NewCore(buildEncoder(config), buildOut(config), config.Level))
		}
	}
	logger := zap.New(zapcore.NewTee(core...), zap.AddCaller())
	zap.ReplaceGlobals(logger)
	log.SetLogger(logger)
}

func buildEncoder(config *LogConfig) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()             // 默认配置
	util.CopyNoneZeroField(&config.EncoderConfig, &encoderConfig) // 覆盖默认配置
	var encoder zapcore.Encoder
	if config.JSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return encoder
}

func buildOut(config *LogConfig) zapcore.WriteSyncer {
	var out []zapcore.WriteSyncer
	for i := range config.Output {
		output := &config.Output[i]
		switch output.Type {
		case LogOutputTypeConsole:
			if output.Destination.Filename == StdErr {
				out = append(out, zapcore.AddSync(os.Stderr))
			} else {
				out = append(out, zapcore.AddSync(os.Stdout))
			}
		case LogOutputTypeFile:
			var fileOut = &lumberjack.Logger{ // 日志切割: 默认配置
				Filename:   AppFilePath("app.log"), // 文件名
				MaxSize:    100,                    // MB 超过这个大小会切割日志
				MaxAge:     30,                     // day 切割的日志最多保存几天
				MaxBackups: 30,                     // 切割的日志最多最多保存几个
				LocalTime:  false,                  // 默认 false=UTC 时间
				Compress:   true,                   // 压缩
			}
			if !filepath.IsAbs(output.Destination.Filename) { // 不是绝对路径 使用相对应用的路径
				// use app dir not working dir(not "./")
				output.Destination.Filename = AppFilePath(output.Destination.Filename)
			}
			util.CopyNoneZeroField(&output.Destination, fileOut) // 覆盖默认配置
			fmt.Println("log  file:", fileOut.Filename)
			out = append(out, zapcore.AddSync(fileOut))
		default:
			out = append(out, zapcore.AddSync(os.Stdout))
		}
	}
	return zapcore.NewMultiWriteSyncer(out...)
}
