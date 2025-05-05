package kjob

import (
	"context"
)

type Config struct {
	Name      string
	startFunc func(ctx context.Context) error
}

func DefaultConfig() *Config {
	return &Config{
		Name:      "",
		startFunc: nil,
	}
}
