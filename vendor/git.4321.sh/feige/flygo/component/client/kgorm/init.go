package kgorm

import (
	"time"

	"git.4321.sh/feige/flygo/core/klog"
	"git.4321.sh/feige/flygo/core/kmetric"
)

func init() {
	type gormStatus struct {
		Gorms map[string]interface{} `json:"gorms"`
	}
	go monitor()
}

func monitor() {
	for {
		time.Sleep(time.Second * 10)
		iterate(func(name string, db *Component) bool {
			sqlDB, err := db.DB()
			if err != nil {
				klog.FlygoLogger.With(klog.FieldComponent(PackageName)).Panic("monitor db error", klog.FieldErr(err))
				return false
			}

			stats := sqlDB.Stats()
			kmetric.LibHandleSummary.Observe(float64(stats.Idle), name, "idle")
			kmetric.LibHandleSummary.Observe(float64(stats.InUse), name, "inuse")
			kmetric.LibHandleSummary.Observe(float64(stats.WaitCount), name, "wait")
			kmetric.LibHandleSummary.Observe(float64(stats.OpenConnections), name, "conns")
			kmetric.LibHandleSummary.Observe(float64(stats.MaxOpenConnections), name, "max_open_conns")
			kmetric.LibHandleSummary.Observe(float64(stats.MaxIdleClosed), name, "max_idle_closed")
			kmetric.LibHandleSummary.Observe(float64(stats.MaxLifetimeClosed), name, "max_lifetime_closed")
			return true
		})
	}
}
