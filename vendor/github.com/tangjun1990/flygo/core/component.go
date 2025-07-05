package core

//Component 组件规范定义
type Component interface {
	ConfigKey() string   // 对应配置文件唯一key
	PackageName() string // 包名
	Start() error        // 启动方法
	Stop() error         // 卸载行为方法
}
