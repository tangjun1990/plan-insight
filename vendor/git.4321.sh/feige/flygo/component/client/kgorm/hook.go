package kgorm

import (
	"fmt"
	"log"
	"time"

	"git.4321.sh/feige/flygo/component/client/kgorm/dsn"

	"git.4321.sh/feige/flygo/core/klog"
	"git.4321.sh/feige/flygo/core/utils/xdebug"
	"gorm.io/gorm"
)

type Handler func(*gorm.DB)

type Processor interface {
	Get(name string) func(*gorm.DB)
	Replace(name string, handler func(*gorm.DB)) error
}

type Hook func(string, *dsn.DSN, string, *config, *klog.Component) func(next Handler) Handler

// debug 拦截器
func debugHook(compName string, dsn *dsn.DSN, op string, options *config, logger *klog.Component) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(db *gorm.DB) {
			beg := time.Now()
			next(db)
			cost := time.Since(beg)

			if db.Error != nil {
				log.Println("[kgorm.response]",
					xdebug.MakeReqResError(compName, fmt.Sprintf("%v", dsn.Addr+"/"+dsn.DBName), cost, logSQL(db.Statement.SQL.String(), db.Statement.Vars, true), db.Error.Error()),
				)
			} else {
				log.Println("[kgorm.response]",
					xdebug.MakeReqResInfo(compName, fmt.Sprintf("%v", dsn.Addr+"/"+dsn.DBName), cost, logSQL(db.Statement.SQL.String(), db.Statement.Vars, true), fmt.Sprintf("%v", db.Statement.Dest)),
				)
			}
		}
	}
}

func logSQL(sql string, args []interface{}, containArgs bool) string {
	if containArgs {
		return bindSQL(sql, args)
	}
	return sql
}
