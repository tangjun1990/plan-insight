package kmetric

import (
	"time"

	"git.4321.sh/feige/flygo/core/kapp"
)

var (
	TypeHTTP         = "http"
	TypeRedis        = "redis"
	TypeGorm         = "gorm"
	TypeWebsocket    = "ws"
	TypeMySQL        = "mysql"
	CodeJobSuccess   = "ok"
	CodeJobFail      = "fail"
	CodeJobReentry   = "reentry"
	CodeCacheMiss    = "miss"
	CodeCacheHit     = "hit"
	DefaultNamespace = "flygo"
)

var (
	ServerHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "server_handle_total",
		Labels:    []string{"type", "method", "peer", "code"},
	}.Build()

	ServerHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "server_handle_seconds",
		Labels:    []string{"type", "method", "peer"},
	}.Build()

	ClientHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_handle_total",
		Labels:    []string{"type", "name", "method", "peer", "code"},
	}.Build()

	ClientHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_handle_seconds",
		Labels:    []string{"type", "name", "method", "peer"},
	}.Build()

	JobHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "job_handle_total",
		Labels:    []string{"type", "name", "code"},
	}.Build()

	JobHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "job_handle_seconds",
		Labels:    []string{"type", "name"},
	}.Build()
	LibHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_seconds",
		Labels:    []string{"type", "method", "address"},
	}.Build()
	LibHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_total",
		Labels:    []string{"type", "method", "address", "code"},
	}.Build()
	LibHandleSummary = SummaryVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_stats",
		Labels:    []string{"name", "status"},
	}.Build()

	CacheHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "cache_handle_total",
		Labels:    []string{"type", "name", "action", "code"},
	}.Build()

	CacheHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "cache_handle_seconds",
		Labels:    []string{"type", "name", "action"},
	}.Build()

	BuildInfoGauge = GaugeVecOpts{
		Namespace: DefaultNamespace,
		Name:      "build_info",
		Labels:    []string{"name", "mode", "region", "zone", "app_version", "flygo_version", "start_time", "build_time", "go_version"},
	}.Build()
)

func init() {
	BuildInfoGauge.WithLabelValues(
		kapp.Name(),
		kapp.AppMode(),
		kapp.AppRegion(),
		kapp.AppZone(),
		kapp.AppVersion(),
		kapp.FlygoVersion(),
		kapp.StartTime(),
		kapp.BuildTime(),
		kapp.GoVersion(),
	).Set(float64(time.Now().UnixNano() / 1e6))
}
