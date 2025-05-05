package kcron

import (
	"sync"
	"time"

	"git.4321.sh/feige/flygo/core/klog"
	"github.com/robfig/cron/v3"
)

func queueIfStillRunning(logger *klog.Component) JobWrapper {
	return func(j Job) Job {
		var mu sync.Mutex
		return cron.FuncJob(func() {
			start := time.Now()
			mu.Lock()
			defer mu.Unlock()
			if dur := time.Since(start); dur > time.Minute {
				logger.Info("cron queue", klog.String("duration", dur.String()))
			}
			j.Run()
		})
	}
}

func skipIfStillRunning(logger *klog.Component) JobWrapper {
	var ch = make(chan struct{}, 1)
	ch <- struct{}{}
	return func(j Job) Job {
		return cron.FuncJob(func() {
			select {
			case v := <-ch:
				j.Run()
				ch <- v
			default:
				logger.Info("cron skip")
			}
		})
	}
}

type Option func(c *Container)

func WithLock(lock Lock) Option {
	return func(c *Container) {
		c.config.lock = lock
	}
}

func WithWrappers(wrappers ...JobWrapper) Option {
	return func(c *Container) {
		if c.config.wrappers == nil {
			c.config.wrappers = []JobWrapper{}
		}
		c.config.wrappers = append(c.config.wrappers, wrappers...)
	}
}

func WithJob(job FuncJob) Option {
	return func(c *Container) {
		c.config.job = job
	}
}

func WithSeconds() Option {
	return func(c *Container) {
		c.config.EnableSeconds = true
	}
}

func WithParser(p cron.Parser) Option {
	return func(c *Container) {
		c.config.parser = p
	}
}

func WithLocation(loc *time.Location) Option {
	return func(c *Container) {
		c.config.loc = loc
	}
}
