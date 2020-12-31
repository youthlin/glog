package common

import (
	"io/ioutil"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

var (
	appConfig      AppConfig // 全局配置
	initConfigOnce sync.Once // 初始化配置一次
)

const (
	LogOutputTypeConsole = "console" // 日志输出到控制台
	LogOutputTypeFile    = "file"    // 日志输出到文件
	StdErr               = "stderr"
)

func mustInitConfig() {
	initConfigOnce.Do(func() {
		file, err := os.Open("./conf/config.yaml")
		if err != nil {
			panic(err)
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(content, &appConfig)
		if err != nil {
			panic(err)
		}
	})
}

func Config() *AppConfig {
	mustInitConfig()
	return &appConfig
}

type AppConfig struct {
	Log []LogConfig `json:"log" yaml:"log"` // 可以同时配置 file, console
	Web WebConfig   `json:"web" yaml:"web"`
}

// 每个配置对应一个 zapcore.Core
type LogConfig struct {
	Name   string          `json:"name" yaml:"name"`     // 起个名字区分，代码中不使用
	Enable bool            `json:"enable" yaml:"enable"` // 是否启用
	JSON   bool            `json:"json" yaml:"json"`     // 是否输出为 json 格式
	Level  zap.AtomicLevel `json:"level" yaml:"level"`   // 日志级别 zapcore.Level: debug, info, warn, error, dpanic, panic, fatal
	Output []LogOutput     `json:"output" yaml:"output"` // 输出目的配置 console/file
	// EncoderConfig 输出格式配置
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
}

type LogOutput struct {
	Type        string            `json:"type" yaml:"type"` // console, file
	Destination lumberjack.Logger `json:"destination" yaml:"destination"`
}
type WebConfig struct {
	Addr string `json:"addr" yaml:"addr"` // :8088
	Mode string `json:"mode" yaml:"mode"` // gin mode: debug, release, test, default=<empty>=debug
}
