package kredis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/tangjun1990/flygo/core/kapp"
	"github.com/tangjun1990/flygo/core/kcfg"
	"github.com/tangjun1990/flygo/core/klog"
)

type Option func(c *Container)

type Container struct {
	config *config
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
	if err := kcfg.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", klog.FieldErr(err), klog.FieldKey(key))
		return c
	}

	c.logger = c.logger.With(klog.FieldComponentName(key))
	c.name = key
	return c
}

func (c *Container) Build(options ...Option) *Component {
	if options == nil {
		options = make([]Option, 0)
	}

	// 全局 debug
	if kapp.IsDevelopmentMode() {
		c.config.Debug = true
	}

	options = append(options, withHook(fixedHook(c.name, c.config, c.logger)))
	if c.config.Debug {
		options = append(options, withHook(debugHook(c.name, c.config, c.logger)))
	}
	for _, option := range options {
		option(c)
	}

	var client redis.Cmdable
	switch c.config.Mode {
	case ClusterMode:
		if len(c.config.Addrs) == 0 {
			c.logger.Panic(`invalid "addrs" config, "addrs" has none addresses but with cluster mode"`)
		}
		client = c.buildCluster()
	case StubMode:
		if c.config.Addr == "" {
			c.logger.Panic(`invalid "addr" config, "addr" is empty but with cluster mode"`)
		}
		client = c.buildStub()
	case SentinelMode:
		if len(c.config.Addrs) == 0 {
			c.logger.Panic(`invalid "addrs" config, "addrs" has none addresses but with sentinel mode"`)
		}
		if c.config.MasterName == "" {
			c.logger.Panic(`invalid "masterName" config, "masterName" is empty but with sentinel mode"`)
		}
		client = c.buildSentinel()
	case RwClusterMode:
		if len(c.config.Addrs) == 0 {
			c.logger.Panic(`invalid "addrs" config, "addrs" has none addresses but with sr mode"`)
		}
		client = c.buildRWCluster()
	default:
		c.logger.Panic(`redis mode must be one of ("stub", "cluster", "sentinel")`)
	}

	c.logger = c.logger.With(klog.FieldAddr(fmt.Sprintf("%s", c.config.Addrs)))
	return &Component{
		config:     c.config,
		client:     client,
		lockClient: &lockClient{client: client},
		logger:     c.logger,
	}
}

func (c *Container) buildCluster() *redis.ClusterClient {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: c.config.Addrs,
		NewClient: func(opt *redis.Options) *redis.Client {
			node := redis.NewClient(opt)

			if c.config.HookLog {
				node.AddHook(NewAccessLogPlugin(opt.Addr, c.config, c.logger))
			}
			return node
		},
		MaxRedirects: c.config.MaxRetries,
		ReadOnly:     c.config.ReadOnly,
		Password:     c.config.Password,
		MaxRetries:   c.config.MaxRetries,
		DialTimeout:  c.config.DialTimeout,
		ReadTimeout:  c.config.ReadTimeout,
		WriteTimeout: c.config.WriteTimeout,
		PoolSize:     c.config.PoolSize,
		MinIdleConns: c.config.MinIdleConns,
		IdleTimeout:  c.config.IdleTimeout,
	})

	for _, incpt := range c.config.hooks {
		clusterClient.AddHook(incpt)
	}

	if err := clusterClient.Ping(context.Background()).Err(); err != nil {
		switch c.config.OnFail {
		case "panic":
			c.logger.Panic("start cluster redis", klog.FieldErr(err))
		default:
			c.logger.Error("start cluster redis", klog.FieldErr(err))
		}
	}
	return clusterClient
}

func (c *Container) buildSentinel() *redis.Client {
	sentinelClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    c.config.MasterName,
		SentinelAddrs: c.config.Addrs,
		Password:      c.config.Password,
		DB:            c.config.DB,
		MaxRetries:    c.config.MaxRetries,
		DialTimeout:   c.config.DialTimeout,
		ReadTimeout:   c.config.ReadTimeout,
		WriteTimeout:  c.config.WriteTimeout,
		PoolSize:      c.config.PoolSize,
		MinIdleConns:  c.config.MinIdleConns,
		IdleTimeout:   c.config.IdleTimeout,
	})

	for _, incpt := range c.config.hooks {
		sentinelClient.AddHook(incpt)
	}

	if err := sentinelClient.Ping(context.Background()).Err(); err != nil {
		switch c.config.OnFail {
		case "panic":
			c.logger.Panic("start sentinel redis", klog.FieldErr(err))
		default:
			c.logger.Error("start sentinel redis", klog.FieldErr(err))
		}
	}
	return sentinelClient
}

func (c *Container) buildStub() *redis.Client {
	stubClient := redis.NewClient(&redis.Options{
		Addr:         c.config.Addr,
		Password:     c.config.Password,
		DB:           c.config.DB,
		MaxRetries:   c.config.MaxRetries,
		DialTimeout:  c.config.DialTimeout,
		ReadTimeout:  c.config.ReadTimeout,
		WriteTimeout: c.config.WriteTimeout,
		PoolSize:     c.config.PoolSize,
		MinIdleConns: c.config.MinIdleConns,
		IdleTimeout:  c.config.IdleTimeout,
	})

	for _, incpt := range c.config.hooks {
		stubClient.AddHook(incpt)
	}

	if c.config.HookLog {
		stubClient.AddHook(NewAccessLogPlugin(c.config.Addr, c.config, c.logger))
	}

	if err := stubClient.Ping(context.Background()).Err(); err != nil {
		switch c.config.OnFail {
		case "panic":
			c.logger.Panic("start stub redis", klog.FieldErr(err))
		default:
			c.logger.Error("start stub redis", klog.FieldErr(err))
		}
	}
	return stubClient
}

func (c *Container) buildRWCluster() *redis.ClusterClient {
	clusterSlots := func(ctx context.Context) ([]redis.ClusterSlot, error) {
		nodes := make([]redis.ClusterNode, len(c.config.Addrs))
		for k, addr := range c.config.Addrs {
			nodes[k] = redis.ClusterNode{
				Addr: addr,
			}
		}
		slots := []redis.ClusterSlot{
			{
				Start: 0,
				End:   16383,
				Nodes: nodes,
			},
		}
		return slots, nil
	}

	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: c.config.Addrs,
		NewClient: func(opt *redis.Options) *redis.Client {
			node := redis.NewClient(opt)

			if c.config.HookLog {
				node.AddHook(NewAccessLogPlugin(opt.Addr, c.config, c.logger))
			}
			return node
		},
		MaxRedirects:  c.config.MaxRetries,
		ReadOnly:      true, // default ReadOnly in this mode
		Password:      c.config.Password,
		MaxRetries:    c.config.MaxRetries,
		DialTimeout:   c.config.DialTimeout,
		ReadTimeout:   c.config.ReadTimeout,
		WriteTimeout:  c.config.WriteTimeout,
		PoolSize:      c.config.PoolSize,
		MinIdleConns:  c.config.MinIdleConns,
		IdleTimeout:   c.config.IdleTimeout,
		ClusterSlots:  clusterSlots,
		RouteRandomly: c.config.RouteRandomly, //默认false，master不承担读请我。如果为true，master也会承接读请求
	})

	for _, h := range c.config.hooks {
		clusterClient.AddHook(h)
	}

	if err := clusterClient.Ping(context.Background()).Err(); err != nil {
		switch c.config.OnFail {
		case "panic":
			c.logger.Panic("start sr redis", klog.FieldErr(err))
		default:
			c.logger.Error("start sr redis", klog.FieldErr(err))
		}
	}
	return clusterClient
}
