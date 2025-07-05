package kgorm

import (
	"github.com/tangjun1990/flygo/component/client/kgorm/dsn"
)

type Option func(c *Container)

func WithDSNParser(parser dsn.DSNParser) Option {
	return func(c *Container) {
		c.dsnParser = parser
	}
}

func WithHook(is ...Hook) Option {
	return func(c *Container) {
		if c.config.hooks == nil {
			c.config.hooks = make([]Hook, 0)
		}
		c.config.hooks = append(c.config.hooks, is...)
	}
}
