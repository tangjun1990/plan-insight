package kcron

import (
	"strings"

	"go.uber.org/zap"

	"github.com/tangjun1990/flygo/core/kcfg"
	"github.com/tangjun1990/flygo/core/klog"

	"github.com/robfig/cron/v3"
)

type Container struct {
	config *Config
	name   string
	logger *klog.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: klog.FlygoLogger.With(klog.FieldComponent(PackageName)),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	if err := kcfg.UnmarshalKey(key, c.config); err != nil {
		c.logger.Panic("parse config error", klog.FieldErr(err), klog.FieldKey(key))
		return c
	}
	c.config.Spec = strings.TrimSpace(c.config.Spec)
	c.logger = c.logger.With(klog.FieldComponentName(key))
	c.name = key
	return c
}

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if c.config.EnableSeconds {
		c.config.parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	}

	switch c.config.DelayExecType {
	case "skip":
		c.config.wrappers = append(c.config.wrappers, skipIfStillRunning(c.logger))
	case "queue":
		c.config.wrappers = append(c.config.wrappers, queueIfStillRunning(c.logger))
	case "concurrent":
	default:
		c.config.wrappers = append(c.config.wrappers, skipIfStillRunning(c.logger))
	}

	if c.config.EnableDistributedTask && c.config.lock == nil {
		c.logger.Panic("lock can not be nil", klog.FieldKey("use WithLock option to set lock"))
	}

	_, err := c.config.parser.Parse(c.config.Spec)
	if err != nil {
		c.logger.Panic("invalid cron spec", zap.Error(err))
	}

	return newComponent(c.name, c.config, c.logger)
}
