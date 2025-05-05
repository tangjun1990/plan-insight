package klog

import (
	"io"

	"git.4321.sh/feige/flygo/core/klog/rotate"
)

func newRotate(config *Config) io.Writer {
	rotatklog := rotate.NewLogger()
	rotatklog.Filename = config.filename()
	rotatklog.MaxSize = config.MaxSize // MB
	rotatklog.MaxAge = config.MaxAge   // days
	rotatklog.MaxBackups = config.MaxBackup
	rotatklog.Interval = config.RotateInterval
	rotatklog.LocalTime = true
	rotatklog.Compress = false
	return rotatklog
}
