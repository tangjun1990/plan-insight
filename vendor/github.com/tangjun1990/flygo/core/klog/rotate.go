package klog

import (
	"io"

	"github.com/tangjun1990/flygo/core/klog/rotate"
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
