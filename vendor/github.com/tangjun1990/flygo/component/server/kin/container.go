package kin

import (
	healthcheck "github.com/RaMin0/gin-health-check"
	"github.com/gin-gonic/gin"
	"github.com/tangjun1990/flygo/core/kapp"
	"github.com/tangjun1990/flygo/core/kcfg"
	"github.com/tangjun1990/flygo/core/klog"
	"github.com/tangjun1990/flygo/core/utils/xnet"
)

// Container 容器
type Container struct {
	config *Config
	name   string
	logger *klog.Component
}

// DefaultContainer 默认容器
func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: klog.FlygoLogger.With(klog.FieldComponent(PackageName)),
	}
}

// Load 加载配置key
func Load(key string) *Container {
	c := DefaultContainer()
	c.logger = c.logger.With(klog.FieldComponentName(key))
	if err := kcfg.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", klog.FieldErr(err), klog.FieldKey(key))
		return c
	}
	var (
		host string
		err  error
	)
	// 获取网卡ip
	if c.config.EnableLocalMainIP {
		host, _, err = xnet.GetLocalMainIP()
		if err != nil {
			host = ""
		}
		c.config.Host = host
	}
	c.name = key
	return c
}

// Build 构建组件
func (c *Container) Build(options ...Option) *Component {
	// 全局 debug
	if kapp.IsDevelopmentMode() {
		c.config.Mode = gin.DebugMode
		c.config.HookReq = true
		c.config.HookRsp = true
	}

	for _, option := range options {
		option(c)
	}
	server := newComponent(c.name, c.config, c.logger)
	server.Use(healthcheck.Default())

	if c.config.HookMetric {
		server.Use(metricServerHook())
	}
	if c.config.HookTrace {
		server.Use(traceServerHook())
	}

	server.Use(defaultServerHook(c.logger, c.config))

	return server
}
