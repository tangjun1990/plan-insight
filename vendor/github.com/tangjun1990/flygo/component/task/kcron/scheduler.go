package kcron

import (
	"sync/atomic"
	"time"

	"github.com/tangjun1990/flygo/core"
)

type immediatelyScheduler struct {
	Schedule
	initOnce uint32
}

func (is *immediatelyScheduler) Next(curr time.Time) (next time.Time) {
	if atomic.CompareAndSwapUint32(&is.initOnce, 0, 1) {
		return curr
	}

	return is.Schedule.Next(curr)
}

type Kcron interface {
	core.Component
}
