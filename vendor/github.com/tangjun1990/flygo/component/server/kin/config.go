package kin

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tangjun1990/flygo/core/kflag"
	"github.com/tangjun1990/flygo/core/utils/xtime"
)

type Config struct {
	Host              string        // 默认0.0.0.0
	Port              int           // 默认8086
	Mode              string        // gin的模式，默认release
	Grace             bool          // 是否开启热更新，默认不开启，热更新基于endless，开启后 kill -1 pid触发热更
	HookMetric        bool          // 默认开启
	HookTrace         bool          // 默认开启
	EnableLocalMainIP bool          // 自动获取ip地址
	SlowLogThreshold  time.Duration // 服务慢日志，默认500ms
	HookReq           bool          // 是否开启记录请求参数，默认不开启
	HookRsp           bool          // 是否开启记录响应参数，默认不开启
	MaxReqContentSize int           // 最大请求 Body 大小，默认 1024
	MaxResContentSize int           // 最大响应结果大小，默认 2048
}

func DefaultConfig() *Config {
	return &Config{
		Host:              kflag.String("host"),
		Port:              8086,
		Mode:              gin.ReleaseMode,
		Grace:             false,
		HookTrace:         true,
		HookMetric:        true,
		SlowLogThreshold:  xtime.Duration("500ms"),
		MaxReqContentSize: 1024,
		MaxResContentSize: 2048,
	}
}

func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
