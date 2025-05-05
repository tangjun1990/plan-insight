package kjob

import (
	"context"
	"git.4321.sh/feige/flygo/core"
	"time"

	"git.4321.sh/feige/flygo/core/kflag"
	"git.4321.sh/feige/flygo/core/klog"
)

func init() {
	kflag.Register(
		&kflag.StringFlag{
			Name:    "job",
			Usage:   "--job",
			Default: "",
		},
	)
}

const PackageName = "core.kjob"

type Component struct {
	name   string
	config *Config
	logger *klog.Component
}

func newComponent(name string, config *Config, logger *klog.Component) *Component {
	return &Component{
		name:   name,
		config: config,
		logger: logger,
	}
}

func (c *Component) ConfigKey() string {
	return c.config.Name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Start() error {
	ctx := context.Background()

	beg := time.Now()
	c.logger.WithCtx(ctx).Info("start kjob", klog.FieldName(c.name))
	err := c.config.startFunc(ctx)
	if err != nil {
		c.logger.WithCtx(ctx).Error("stop kjob", klog.FieldName(c.name), klog.FieldErr(err), klog.FieldDuration(time.Since(beg)))
	} else {
		c.logger.WithCtx(ctx).Info("stop kjob", klog.FieldName(c.name), klog.FieldDuration(time.Since(beg)))
	}
	return err
}

func (c *Component) Stop() error {
	return nil
}

type Kjob interface {
	core.Component
}
