package kgorm

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tangjun1990/flygo/core/ktrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/tangjun1990/flygo/component/client/kgorm/dsn"

	"github.com/tangjun1990/flygo/core/klog"
	"github.com/tangjun1990/flygo/core/utils/xdebug"
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

func peerInfo(addr string) (hostname string, port int) {
	if idx := strings.IndexByte(addr, ':'); idx >= 0 {
		hostname = addr[:idx]
		port, _ = strconv.Atoi(addr[idx+1:])
	}
	return hostname, port
}

func traceHook(compName string, dsn *dsn.DSN, op string, options *config, logger *klog.Component) func(Handler) Handler {
	ip, port := peerInfo(dsn.Addr)
	attrs := []attribute.KeyValue{
		semconv.NetHostIPKey.String(ip),
		semconv.NetPeerPortKey.Int(port),
		semconv.NetTransportKey.String(dsn.Net),
		semconv.DBNameKey.String(dsn.DBName),
	}
	tracer := ktrace.NewTracer(trace.SpanKindClient)
	return func(next Handler) Handler {
		return func(db *gorm.DB) {
			if db.Statement.Context != nil {
				operation := "gorm:"
				if len(db.Statement.BuildClauses) > 0 {
					operation += strings.ToLower(db.Statement.BuildClauses[0])
				}
				if db.Statement.Table != "" {
					operation += ":" + db.Statement.Table
				}

				_, span := tracer.Start(db.Statement.Context, operation, nil, trace.WithAttributes(attrs...))
				defer span.End()
				// 延迟执行 scope.CombinedConditionSql() 避免sqlVar被重复追加
				next(db)
				span.SetAttributes(
					semconv.DBSystemKey.String(db.Dialector.Name()),
					semconv.DBStatementKey.String(logSQL(db.Statement.SQL.String(), db.Statement.Vars, true)),
					semconv.DBOperationKey.String(operation),
					semconv.DBSQLTableKey.String(db.Statement.Table),
					semconv.NetPeerNameKey.String(dsn.Addr),
					attribute.Int64("db.rows_affected", db.RowsAffected),
				)
				if db.Error != nil {
					span.RecordError(db.Error)
					span.SetStatus(codes.Error, db.Error.Error())
					return
				}
				span.SetStatus(codes.Ok, "OK")
				return
			}

			next(db)
		}
	}
}
