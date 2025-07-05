package kgorm

import (
	"sync"

	"github.com/tangjun1990/flygo/core/klog"
)

var instances = sync.Map{}

func iterate(fn func(name string, db *Component) bool) {
	instances.Range(func(key, val interface{}) bool {
		return fn(key.(string), val.(*Component))
	})
}

func configs() map[string]interface{} {
	var rets = make(map[string]interface{})
	instances.Range(func(key, val interface{}) bool {
		return true
	})

	return rets
}

func stats() (stats map[string]interface{}) {
	stats = make(map[string]interface{})
	instances.Range(func(key, val interface{}) bool {
		name := key.(string)
		db := val.(*Component)

		sqlDB, err := db.DB()
		if err != nil {
			klog.FlygoLogger.With(klog.FieldComponent(PackageName)).Panic("stats db error", klog.FieldErr(err))
			return false
		}
		stats[name] = sqlDB.Stats()
		return true
	})

	return
}
