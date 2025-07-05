package kredis

import (
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/tangjun1990/flygo/core/utils/xtime"
)

const (
	// ClusterMode using clusterClient
	ClusterMode string = "cluster"
	// StubMode using stubClient
	StubMode string = "stub"
	// SentinelMode using Failover sentinel client
	SentinelMode string = "sentinel"
	// RwClusterMode without using cluster mode or Redis Sentinel.
	RwClusterMode string = "rwCluster"
)

// config for redis, contains RedisStubConfig, RedisClusterConfig and RedisSentinelConfig
type config struct {
	Addrs            []string      // Addrs 实例配置地址
	Addr             string        // Addr stubConfig 实例配置地址
	Mode             string        // Mode Redis模式 cluster|rwCluster|stub|sentinel
	MasterName       string        // MasterName 哨兵主节点名称，sentinel模式下需要配置此项
	Password         string        // Password 密码
	DB               int           // DB，默认为0, 一般应用不推荐使用DB分片
	PoolSize         int           // PoolSize 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	MaxRetries       int           // MaxRetries 网络相关的错误最大重试次数 默认8次
	MinIdleConns     int           // MinIdleConns 最小空闲连接数
	DialTimeout      time.Duration // DialTimeout 拨超时时间
	ReadTimeout      time.Duration // ReadTimeout 读超时 默认3s
	WriteTimeout     time.Duration // WriteTimeout 读超时 默认3s
	IdleTimeout      time.Duration // IdleTimeout 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	RouteRandomly    bool          // rwCluster模式下有效，如果开启，那么master也会分担读请求
	Debug            bool          // Debug开关
	ReadOnly         bool          // ReadOnly 集群模式 在从属节点上启用读模式
	SlowLogThreshold time.Duration // 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	OnFail           string        // OnFail panic|error

	// 日志开关
	HookLog bool // 是否开启，记录请求数据
	HookReq bool // 是否开启记录请求参数
	HookRsp bool // 是否开启记录响应参数

	MaxResContentSize int //最大响应结果大小

	hooks []redis.Hook // redis 拦截器
}

func DefaultConfig() *config {
	return &config{
		Mode:              StubMode,
		DB:                0,
		PoolSize:          10,
		MaxRetries:        3,
		MinIdleConns:      100,
		DialTimeout:       xtime.Duration("1s"),
		ReadTimeout:       xtime.Duration("1s"),
		WriteTimeout:      xtime.Duration("1s"),
		IdleTimeout:       xtime.Duration("60s"),
		RouteRandomly:     false,
		ReadOnly:          false,
		Debug:             false,
		SlowLogThreshold:  xtime.Duration("250ms"),
		OnFail:            "panic",
		MaxResContentSize: 2048,
	}
}
