package kcron

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"go.uber.org/zap"

	"github.com/tangjun1990/flygo/core/klog"
	"github.com/tangjun1990/flygo/core/utils/xstring"
)

const PackageName = "core.kcron"

type (
	JobWrapper = cron.JobWrapper
	EntryID    = cron.EntryID
	Entry      = cron.Entry
	Schedule   = cron.Schedule
	Parser     = cron.Parser

	Job      = cron.Job
	NamedJob interface {
		Run(ctx context.Context) error
		Name() string
	}
)

type FuncJob func(ctx context.Context) error

func (f FuncJob) Run(ctx context.Context) error {
	return f(ctx)
}

func (f FuncJob) Name() string { return xstring.FunctionName(f) }

type Component struct {
	name   string
	config *Config
	cron   *cron.Cron
	logger *klog.Component
}

func newComponent(name string, config *Config, logger *klog.Component) *Component {
	return &Component{
		config: config,
		cron: cron.New(
			cron.WithParser(config.parser),
			cron.WithChain(config.wrappers...),
			cron.WithLogger(&wrappedLogger{logger}),
			cron.WithLocation(config.loc),
		),
		name:   name,
		logger: logger,
	}
}

func (c *Component) ConfigKey() string {
	return c.name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Start() error {
	if !c.config.Enable {
		return nil
	}

	if c.config.EnableDistributedTask {
		go c.startDistributedTask()
	} else {
		err := c.startTask()
		if err != nil {
			return err
		}
	}

	c.cron.Run()
	return nil
}

func (c *Component) Stop() error {
	_ = c.cron.Stop()
	if c.config.EnableDistributedTask {
		ctx, cancel := context.WithTimeout(context.Background(), c.config.WaitUnlockTime)
		defer cancel()
		err := c.config.lock.Unlock(ctx)
		if err != nil {
			c.logger.WithCtx(ctx).Info("mutex unlock", klog.FieldErr(err))
			return fmt.Errorf("cron stop err: %w", err)
		}
	}
	return nil
}

func (c *Component) schedule(schedule Schedule, job NamedJob) EntryID {
	if c.config.EnableImmediatelyRun {
		schedule = &immediatelyScheduler{
			Schedule: schedule,
		}
	}
	innerJob := &wrappedJob{
		NamedJob: job,
		logger:   c.logger,
	}
	c.logger.WithCtx(context.Background()).Info("add job", klog.String("name", job.Name()))
	return c.cron.Schedule(schedule, innerJob)
}

func (c *Component) addJob(spec string, cmd NamedJob) (EntryID, error) {
	schedule, err := c.config.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.schedule(schedule, cmd), nil
}

func (c *Component) removeJob(id EntryID) {
	c.cron.Remove(id)
}

func (c *Component) startDistributedTask() {
	for {
		func() {
			defer time.Sleep(c.config.RefreshGap)
			ctx, cancel := context.WithTimeout(context.Background(), c.config.WaitLockTime)

			err := c.config.lock.Lock(ctx, c.config.LockTTL)
			cancel()
			if err != nil {
				c.logger.WithCtx(ctx).Info("job lock not obtained", klog.FieldErr(err))
				return
			}

			c.logger.WithCtx(ctx).Info("add cron", klog.Int("number of scheduled jobs", len(c.cron.Entries())))

			entryID, err := c.addJob(c.config.Spec, c.config.job)
			if err != nil {
				c.logger.WithCtx(ctx).Error("add job failed", zap.Error(err))
				return
			}

			err = c.keepLockAlive()
			if err != nil {
				c.logger.WithCtx(ctx).Error("job lost", zap.String("name", c.name), zap.Error(err))
			}

			c.removeJob(entryID)
		}()
	}
}

func (c *Component) keepLockAlive() error {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), c.config.WaitLockTime)
		err := c.config.lock.Refresh(ctx, c.config.LockTTL)
		cancel()
		if err != nil {
			c.logger.WithCtx(ctx).Info("mutex lock", klog.FieldErr(err))
			return err
		}

		time.Sleep(c.config.RefreshGap)
	}
}

func (c *Component) startTask() (err error) {
	_, err = c.addJob(c.config.Spec, c.config.job)
	if err != nil {
		return
	}

	c.logger.WithCtx(context.Background()).Info("add cron", klog.Int("number of scheduled jobs", len(c.cron.Entries())))
	return nil
}
