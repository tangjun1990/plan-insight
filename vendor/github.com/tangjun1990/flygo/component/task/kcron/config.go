package kcron

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tangjun1990/flygo/core/utils/xtime"
)

type Config struct {
	Spec string

	WaitLockTime   time.Duration
	LockTTL        time.Duration
	RefreshGap     time.Duration
	WaitUnlockTime time.Duration

	DelayExecType         string
	Enable                bool
	EnableDistributedTask bool
	EnableImmediatelyRun  bool
	EnableSeconds         bool

	wrappers []JobWrapper
	parser   cron.Parser
	lock     Lock
	job      FuncJob
	loc      *time.Location
}

func DefaultConfig() *Config {
	return &Config{
		Spec:                  "",
		WaitLockTime:          xtime.Duration("4s"),
		LockTTL:               xtime.Duration("16s"),
		RefreshGap:            xtime.Duration("4s"),
		WaitUnlockTime:        xtime.Duration("1s"),
		DelayExecType:         "skip",
		Enable:                true,
		EnableDistributedTask: false,
		EnableImmediatelyRun:  false,
		EnableSeconds:         false,
		wrappers:              []JobWrapper{},
		parser:                cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		lock:                  nil,
		job:                   nil,
		loc:                   time.Local,
	}
}
