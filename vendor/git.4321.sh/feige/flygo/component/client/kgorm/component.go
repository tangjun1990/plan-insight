package kgorm

import (
	"context"
	"errors"
	"git.4321.sh/feige/flygo/component/client/kgorm/dsn"
	"git.4321.sh/feige/flygo/core/klog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

const PackageName = "component.kgorm"

var (
	errSlowCommand        = errors.New("mysql slow command")
	ErrRecordNotFound     = gorm.ErrRecordNotFound
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
)

type (
	DB             gorm.DB
	Dialector      = gorm.Dialector
	Model          = gorm.Model
	Field          = schema.Field
	Association    = gorm.Association
	NamingStrategy = schema.NamingStrategy
	Logger         = logger.Interface
)

type Component = gorm.DB

func WithContext(c context.Context, db *Component) *Component {
	db.Statement.Context = c
	return db
}

func WithCtx(c context.Context, db *Component) *Component {
	return WithContext(c, db)
}

func newComponent(compName string, dsnParser dsn.DSNParser, config *config, klogger *klog.Component) (*Component, error) {
	db, err := gorm.Open(dsnParser.GetDialector(config.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if config.RawDebug {
		db = db.Debug()
	}

	gormDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 创建 从库连接
	err = useSlavesDB(db, dsnParser, config)
	if err != nil {
		return nil, err
	}

	gormDB.SetMaxIdleConns(config.MaxIdleConns)
	gormDB.SetMaxOpenConns(config.MaxOpenConns)

	if config.ConnMaxLifetime != 0 {
		gormDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	}

	replace := func(processor Processor, callbackName string, hooks ...Hook) {
		handler := processor.Get(callbackName)
		for _, hook := range config.hooks {
			handler = hook(compName, config.dsnCfg, callbackName, config, klogger)(handler)
		}
		processor.Replace(callbackName, handler)
	}

	replace(db.Callback().Create(), "gorm:create", config.hooks...)
	replace(db.Callback().Update(), "gorm:update", config.hooks...)
	replace(db.Callback().Delete(), "gorm:delete", config.hooks...)
	replace(db.Callback().Query(), "gorm:query", config.hooks...)
	replace(db.Callback().Raw(), "gorm:raw", config.hooks...)

	return db, nil
}

// 创建从库连接
func useSlavesDB(db *gorm.DB, dsnParser dsn.DSNParser, c *config) error {
	if len(c.SlavesDSN) == 0 {
		return nil
	}
	list := make([]gorm.Dialector, 0)
	for _, itemDSN := range c.SlavesDSN {
		list = append(list, dsnParser.GetDialector(itemDSN))
	}
	resolver := dbresolver.Register(dbresolver.Config{
		Replicas: list,
	})

	resolver.SetMaxIdleConns(c.MaxIdleConns)
	resolver.SetMaxOpenConns(c.MaxOpenConns)
	if c.ConnMaxLifetime != 0 {
		resolver.SetConnMaxLifetime(c.ConnMaxLifetime)
	}
	return db.Use(resolver)
}
