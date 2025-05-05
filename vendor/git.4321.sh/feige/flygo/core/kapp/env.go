package kapp

import (
	"os"
	"strings"
)

const (
	// EnvAppName 应用名环境变量
	EnvAppName = "APP_NAME"
	// EnvAppMode 应用模式环境变量
	EnvAppMode = "APP_MODE"
	// EnvAppRegion ...
	EnvAppRegion = "APP_REGION"
	// EnvAppZone ...
	EnvAppZone = "APP_ZONE"
	// EnvAppHost ...
	EnvAppHost = "APP_HOST"
	// EnvAppInstance 应用实例ID环境变量
	EnvAppInstance = "APP_INSTANCE"
	// FlygoDebug 调试环境变量，export FLYGO_DEBUG=1，开启应用的调试模式
	FlygoDebug = "FLYGO_DEBUG"
	// AppConfigPath 应用配置环境变量
	AppConfigPath = "APP_CONFIG_PATH"
	// AppLogPath 应用日志环境变量
	AppLogPath = "APP_LOG_PATH"
	// AppLogAddApp 应用日志增加应用名环境变量，如果增加该环境变量，日志里会将应用名写入到app字段里
	AppLogAddApp    = "APP_LOG_ADD_APP"
	AppLogExtraKeys = "APP_LOG_EXTRA_KEYS"
	// AppTraceIDName 应用链路ID环境变量，不配置，默认x-trace-id
	AppTraceIDName = "APP_TRACE_ID_NAME"
	// DefaultDeployment ...
	DefaultDeployment = ""
)

var (
	appMode        string
	appRegion      string
	appZone        string
	appInstance    string
	okDebug        string
	okLogPath      string
	okTraceIDName  string
	okLogExtraKeys []string
)

func initEnv() {
	appMode = os.Getenv(EnvAppMode)
	appRegion = os.Getenv(EnvAppRegion)
	appZone = os.Getenv(EnvAppZone)
	appInstance = os.Getenv(EnvAppInstance)
	if appInstance == "" {
		appInstance = HostName()
	}
	okDebug = os.Getenv(FlygoDebug)
	okLogPath = os.Getenv(AppLogPath)
	okTraceIDName = os.Getenv(AppTraceIDName)
	if okTraceIDName == "" {
		okTraceIDName = "x-trace-id"
	}
	okLogExtraKeys = strings.Split(os.Getenv(AppLogExtraKeys), ",")
}

func AppMode() string {
	return appMode
}

func AppRegion() string {
	return appRegion
}

func AppZone() string {
	return appZone
}

func AppInstance() string {
	return appInstance
}

func IsDevelopmentMode() bool {
	return okDebug == "1"
}

func OkLogPath() string {
	return okLogPath
}

func OkLogExtraKeys() []string {
	return okLogExtraKeys
}
