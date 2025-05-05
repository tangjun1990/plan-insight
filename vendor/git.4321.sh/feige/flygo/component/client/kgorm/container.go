package kgorm

import (
	dsn2 "git.4321.sh/feige/flygo/component/client/kgorm/dsn"
	"git.4321.sh/feige/flygo/core/kapp"
	"git.4321.sh/feige/flygo/core/kcfg"
	"git.4321.sh/feige/flygo/core/klog"
	"git.4321.sh/feige/flygo/core/kmetric"
)

// Container ...
type Container struct {
	name      string
	config    *config
	logger    *klog.Component
	dsnParser dsn2.DSNParser
}

// DefaultContainer ...
func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: klog.FlygoLogger.With(klog.FieldComponent(PackageName)),
	}
}

// Load ...
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

func (c *Container) setDSNParserIfNotExists(dialect string) error {
	if c.dsnParser != nil {
		return nil
	}
	switch dialect {
	case dialectMysql:
		c.dsnParser = dsn2.DefaultMysqlDSNParser
	case dialectClickhouse:
		c.dsnParser = dsn2.DefaultClickhouseDSNParser
	default:
		return errSupportDialect
	}
	return nil
}

// Build 构建组件
func (c *Container) Build(options ...Option) *Component {

	// 全局 debug
	if kapp.IsDevelopmentMode() {
		c.config.Debug = true
	}

	// 加载拦截器
	if c.config.Debug {
		options = append(options, WithHook(debugHook))
	}

	for _, option := range options {
		option(c)
	}

	// 加载 dsn 解析器
	var err error
	err = c.setDSNParserIfNotExists(c.config.Dialect)
	if err != nil {
		c.logger.Panic("setDSNParserIfNotExists err", klog.String("dialect", c.config.Dialect), klog.FieldErr(err))
	}
	// 解析 dsn
	c.config.dsnCfg, err = c.dsnParser.ParseDSN(c.config.DSN)

	fields := make([]klog.Field, 0)

	fields = append(fields, klog.FieldNetPeerIp(c.config.dsnCfg.Addr))

	if err == nil {
		c.logger.With(fields...).Info("start db", klog.FieldName(c.config.dsnCfg.DBName))
	} else {
		c.logger.With(fields...).Panic("start db", klog.FieldErr(err))
	}
	// 设置 logger

	c.logger = c.logger.With(fields...)

	// 初始化组件
	component, err := newComponent(c.name, c.dsnParser, c.config, c.logger)
	if err != nil {
		if c.config.OnFail == "panic" {
			c.logger.With(fields...).Panic("open db", klog.FieldErrKind("register err"), klog.FieldErr(err), klog.FieldValueAny(c.config))
		} else {
			kmetric.ClientHandleCounter.Inc(kmetric.TypeGorm, c.name, c.name+".ping", c.config.dsnCfg.Addr, "open err")
			c.logger.With(fields...).Error("open db", klog.FieldErrKind("register err"), klog.FieldErr(err), klog.FieldValueAny(c.config))
			return component
		}
	}

	sqlDB, err := component.DB()
	if err != nil {
		c.logger.With(fields...).Panic("ping db", klog.FieldErrKind("register err"), klog.FieldErr(err), klog.FieldValueAny(c.config))
	}
	if err := sqlDB.Ping(); err != nil {
		c.logger.With(fields...).Panic("ping db", klog.FieldErrKind("register err"), klog.FieldErr(err), klog.FieldValueAny(c.config))
	}

	instances.Store(c.name, component)
	return component
}
