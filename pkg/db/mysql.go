package db

import (
	"fmt"
	"sync"

	"github.com/tangjun1990/flygo/component/client/kgorm"
	"github.com/tangjun1990/flygo/core/kcfg"
)

var _dbs = make(map[string]*kgorm.Component)
var _dbMu sync.RWMutex
var _configKey = "mysql"

// InitBatch 批量初始化
func InitBatch() error {
	configMap := kcfg.GetStringMap(_configKey)
	if len(configMap) == 0 {
		panic("no find mysql config")
	}
	for dbName := range configMap {
		initOne(dbName)
	}
	return nil
}

func initOne(name string) *kgorm.Component {
	_dbMu.Lock()
	defer _dbMu.Unlock()
	_dbs[name] = kgorm.Load(_configKey + "." + name).Build()
	return _dbs[name]
}

// GetDB 根据名称获取数据库链接
func GetDB(name string) (*kgorm.Component, bool) {
	_dbMu.RLock()
	db, ok := _dbs[name]
	_dbMu.RUnlock()
	if !ok {
		return initOne(name), true
	}
	return db, ok
}

// GetDBMust 根据名称必须获取数据库链接
func GetDBMust(name string) *kgorm.Component {
	db, ok := GetDB(name)
	if !ok {
		panic(fmt.Sprintf("get database fail, name:%s", name))
	}
	return db
}
