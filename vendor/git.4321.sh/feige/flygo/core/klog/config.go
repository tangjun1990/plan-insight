package klog

import (
	"fmt"
	"time"

	"git.4321.sh/feige/flygo/core/kapp"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config ...
type Config struct {
	Debug               bool          // 是否双写至文件控制日志输出到终端
	Level               string        // 日志初始等级，默认info级别
	Dir                 string        // [fileWriter]日志输出目录，默认logs
	Name                string        // [fileWriter]日志文件名称，默认框架日志flygo.sys，业务日志default.log
	MaxSize             int           // [fileWriter]日志输出文件最大长度，超过改值则截断，默认500M
	MaxAge              int           // [fileWriter]日志存储最大时间，默认最大保存天数为7天
	MaxBackup           int           // [fileWriter]日志存储最大数量，默认最大保存文件个数为10个
	RotateInterval      time.Duration // [fileWriter]日志轮转时间，默认1天
	EnableAddCaller     bool          // 是否添加调用者信息，默认不加调用者信息
	EnableAsync         bool          // 是否异步，默认异步
	FlushBufferSize     int           // 缓冲大小，默认256 * 1024B
	FlushBufferInterval time.Duration // 缓冲时间，默认5秒
	Writer              string        // 使用哪种Writer，可选[file|stderr]，默认file

	fields        []zap.Field // 日志初始化字段
	CallerSkip    int
	encoderConfig *zapcore.EncoderConfig
}

const (
	writerRotateFile = "file"
	writerStderr     = "stderr"
)

// filename ...
func (config *Config) filename() string {
	return fmt.Sprintf("%s/%s", config.Dir, config.Name)
}

// DefaultConfig ...
func DefaultConfig() *Config {
	dir := "./logs"
	if kapp.OkLogPath() != "" {
		dir = kapp.OkLogPath()
	}
	return &Config{
		Name:                DefaultLoggerName,
		Dir:                 dir,
		Level:               "info",
		FlushBufferSize:     defaultBufferSize,
		FlushBufferInterval: defaultFlushInterval,
		MaxSize:             500, // 500M
		MaxAge:              7,   // 1 day
		MaxBackup:           10,  // 10 backup
		RotateInterval:      24 * time.Hour,
		CallerSkip:          1,
		EnableAddCaller:     true,
		EnableAsync:         true,
		encoderConfig:       defaultZapConfig(),
		Writer:              writerRotateFile,
	}
}
