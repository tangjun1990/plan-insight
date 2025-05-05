package kcron

import (
	"git.4321.sh/feige/flygo/core"
	"sync/atomic"
	"time"
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
