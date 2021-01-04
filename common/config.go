package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

var (
	appConfig      AppConfig // 全局配置
	initConfigOnce sync.Once // 初始化配置一次
	appDir         = "."     // 应用路径
)

const (
	LogOutputTypeConsole = "console" // 日志输出到控制台
	LogOutputTypeFile    = "file"    // 日志输出到文件
	StdErr               = "stderr"
)

func mustInitConfig() {
	initConfigOnce.Do(func() {
		executable, err := os.Executable() // 可执行文件完整路径
		if err != nil {
			executable = os.Args[0]
		}
		executable, err = filepath.EvalSymlinks(executable) // 如果是符号链接 解析出实际路径
		if err != nil {
			executable = os.Args[0]
		}
		appDir = filepath.Dir(executable)           // 获得可执行文件所在路径
		confFile := AppFilePath("conf/config.yaml") // 配置文件完整路径
		if _, err = os.Stat(confFile); os.IsNotExist(err) {
			appDir = "." // run on GoLand 不存在配置文件，可能是在 GoLand 里直接启动的，使用工作路径(PWD)
		}
		appDir, err = filepath.Abs(appDir) // 绝对路径
		if err != nil {
			appDir = "."
		}
		fmt.Println("app  dir:", appDir)
		confFile = AppFilePath("conf/config.yaml")
		fmt.Println("conf file:", confFile)
		file, err := os.Open(confFile)
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

// Config return AppConfig
func Config() *AppConfig {
	mustInitConfig()
	return &appConfig
}

// AppDir return app dir
func AppDir() string {
	mustInitConfig()
	return appDir
}

// AppFilePath return file in the app dir
func AppFilePath(name string) string {
	join := filepath.Join(appDir, name)
	abs, err := filepath.Abs(join)
	if err != nil {
		return join
	}
	return abs
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
