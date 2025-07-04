package redis

import (
	"sync"

	"github.com/tangjun1990/flygo/component/client/kredis"
	"github.com/tangjun1990/flygo/core/kcfg"
)

var _redisMap = make(map[string]*kredis.Component)
var _redisMu sync.RWMutex
var _configKey = "redis"

// InitBatch 批量初始化redis
func InitBatch() error {
	configMap := kcfg.GetStringMap(_configKey)
	if len(configMap) == 0 {
		panic("no find redis config")
	}
	for name := range configMap {
		initOne(name)
	}
	return nil
}

func initOne(name string) {
	_redisMu.Lock()
	defer _redisMu.Unlock()
	_redisMap[name] = kredis.Load(_configKey + "." + name).Build()
}

// GetRedis 通过名称获取redis实例
func GetRedis(name string) (*kredis.Component, bool) {
	_redisMu.RLock()
	defer _redisMu.RUnlock()
	ins, ok := _redisMap[name]
	return ins, ok
}

// DefaultRedis 获取默认的redis实例
func DefaultRedis() (*kredis.Component, bool) {
	return GetRedis("default")
}

// GetRedisMust 通过链接名称获取redis实例，获取失败则报错
func GetRedisMust(name string) *kredis.Component {
	_redisMu.RLock()
	defer _redisMu.RUnlock()
	ins, ok := _redisMap[name]
	if !ok {
		panic("no redis instance named:" + name)
	}
	return ins
}

// DefaultRedisMust 获取默认的 redis 实例
func DefaultRedisMust() *kredis.Component {
	return GetRedisMust("default")
}
