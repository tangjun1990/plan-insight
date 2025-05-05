package klog

import (
	"git.4321.sh/feige/flygo/core/kapp"
	"git.4321.sh/feige/flygo/core/kcfg"
	"git.4321.sh/feige/flygo/core/utils/xnet"
	"os"
)

// Container 容器
type Container struct {
	config *Config
	name   string
}

// DefaultContainer 默认容器
func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
	}
}

// Load 加载配置key
func Load(key string) *Container {
	c := DefaultContainer()
	if err := kcfg.UnmarshalKey(key, &c.config); err != nil {
		panic(err)
	}
	c.name = key
	return c
}

// Build 构建组件
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	// 全局 debug
	if kapp.IsDevelopmentMode() {
		c.config.Debug = true           // 调试模式，终端输出
		c.config.EnableAsync = false    // 调试模式，同步输出
		c.config.EnableAddCaller = true // 调试模式，增加行号输出
	}

	if c.config.encoderConfig == nil {
		c.config.encoderConfig = defaultZapConfig()
	}

	if c.config.Debug {
		c.config.encoderConfig = defaultDebugConfig()
	}

	// 添加日志基础字段
	c.baseField()

	logger := newLogger(c.name, c.config)
	if c.name != "" {
		logger.AutoLevel(c.name + ".level")
	}

	return logger
}

// 日志基础字段
func (c *Container) baseField() {
	serverIP, _, _ := xnet.GetLocalMainIP()
	c.config.fields = append(c.config.fields,
		// service_name 服务名称
		FieldServiceName(kapp.ServiceName()),
		// server_ip 服务端 IP
		FieldNetHostIp(serverIP),
		FieldDeploymentEnvironment(kapp.Environment()),
	)

	if os.Getenv("NODEIP") != "" {
		c.config.fields = append(c.config.fields,
			FieldK8sNodeIp(os.Getenv("NODEIP")),
		)
	}
	if os.Getenv("IDC") != "" {
		c.config.fields = append(c.config.fields,
			FieldCloudRegion(os.Getenv("IDC")),
		)
	}

}
