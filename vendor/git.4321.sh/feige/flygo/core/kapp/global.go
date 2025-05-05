package kapp

import (
	"errors"
)

const (
	GlobalKey = "app.global" // 全局配置
)

var (
	GlobalConfig = &Global{
		Environment: "staging",
	}

	errorServiceNameEmpty = errors.New("serviceName is required")
)

type Global struct {
	ServiceName string
	Environment string
}

// CheckAndInit 验证
func (g *Global) CheckAndInit() error {
	if g.ServiceName == "" {
		return errorServiceNameEmpty
	}

	SetAppName(g.ServiceName)
	return nil
}

func ServiceName() string {
	return GlobalConfig.ServiceName
}

func Environment() string {
	return GlobalConfig.Environment
}
