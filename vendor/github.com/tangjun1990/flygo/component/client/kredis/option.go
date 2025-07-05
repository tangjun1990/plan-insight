package kredis

import (
	"github.com/tangjun1990/flygo/core/klog"

	"github.com/go-redis/redis/v8"
)

func WithStub() Option {
	return func(c *Container) {
		c.config.Mode = StubMode
	}
}

func WithCluster() Option {
	return func(c *Container) {
		c.config.Mode = ClusterMode
	}
}

func WithSentinel() Option {
	return func(c *Container) {
		c.config.Mode = SentinelMode
	}
}

func withHook(hooks ...redis.Hook) Option {
	return func(c *Container) {
		if c.config.hooks == nil {
			c.config.hooks = make([]redis.Hook, 0, len(hooks))
		}
		c.config.hooks = append(c.config.hooks, hooks...)
	}
}

func WithPassword(password string) Option {
	return func(c *Container) {
		c.config.Password = password
	}
}

func WithAddr(addr string) Option {
	return func(c *Container) {
		c.config.Addr = addr
	}
}

func WithAddrs(addrs []string) Option {
	return func(c *Container) {
		c.config.Addrs = addrs
	}
}

// 设置日志组件
func WithLogger(logger *klog.Component) Option {
	return func(c *Container) {
		c.logger = logger
	}
}
