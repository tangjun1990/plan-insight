package kjob

import "github.com/tangjun1990/flygo/core/klog"

type Container struct {
	config *Config
	logger *klog.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: klog.FlygoLogger.With(klog.FieldComponent(PackageName)),
	}
}

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.config.Name, c.config, c.logger)
}
