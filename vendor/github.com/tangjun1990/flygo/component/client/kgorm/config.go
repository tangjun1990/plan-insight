package kgorm

import (
	"errors"
	"time"

	"github.com/tangjun1990/flygo/component/client/kgorm/dsn"

	"github.com/tangjun1990/flygo/core/utils/xtime"
)

const (
	dialectMysql      = "mysql"
	dialectClickhouse = "clickhouse"
)

var errSupportDialect = errors.New("invalid support Dialect")

type config struct {
	Dialect          string        // 选择数据库种类，默认mysql
	DSN              string        // DSN地址: mysql://root:secret@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=1s&readTimeout=3s&writeTimeout=3s
	SlavesDSN        []string      // 从库 DSN 地址
	Debug            bool          // 是否开启调试，默认不开启，开启后并加上export FLYGO=true，可以看到每次请求，配置名、地址、耗时、请求数据、响应数据
	RawDebug         bool          // 是否开启原生调试开关，默认不开启
	MaxIdleConns     int           // 最大空闲连接数，默认10
	MaxOpenConns     int           // 最大活动连接数，默认100
	ConnMaxLifetime  time.Duration // 连接的最大存活时间，默认300s
	OnFail           string        // 创建连接的错误级别，=panic时，如果创建失败，立即panic，默认连接不上panic
	SlowLogThreshold time.Duration // 慢日志阈值，默认500ms
	EnableDetailSQL  bool          // 记录错误sql时,是否打印包含参数的完整sql语句，select * from aid = ?;

	HookTrace         bool // 是否开启链路
	HookReq           bool // 是否开启，记录请求数据
	HookRsp           bool // 是否开启记录响应参数
	MaxResContentSize int  // 最大响应结果大小
	hooks             []Hook
	dsnCfg            *dsn.DSN
}

func DefaultConfig() *config {
	return &config{
		DSN:               "",
		Dialect:           "mysql",
		Debug:             false,
		MaxIdleConns:      10,
		MaxOpenConns:      100,
		ConnMaxLifetime:   xtime.Duration("300s"),
		OnFail:            "panic",
		SlowLogThreshold:  xtime.Duration("500ms"),
		MaxResContentSize: 2048,
	}
}
