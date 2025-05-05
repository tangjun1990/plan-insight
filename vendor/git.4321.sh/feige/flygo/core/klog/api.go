package klog

import (
	"context"
)

const PackageName = "core.klog"

var DefaultLogger *Component
var FlygoLogger *Component

func init() {
	DefaultLogger = DefaultContainer().Build(WithFileName(DefaultLoggerName))
	FlygoLogger = DefaultContainer().Build(WithFileName(FlygoLoggerName))
}

const (
	DefaultLoggerName = "default.log"
	FlygoLoggerName   = "flygo.sys"
)

func getDefaultLogger() *Component {
	return DefaultLogger.WithCallerSkip(1)
}

func Info(msg string, fields ...Field) {
	getDefaultLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...Field) {
	getDefaultLogger().Debug(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	getDefaultLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	getDefaultLogger().Error(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	getDefaultLogger().Panic(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	getDefaultLogger().DPanic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	getDefaultLogger().Fatal(msg, fields...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	getDefaultLogger().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	getDefaultLogger().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	getDefaultLogger().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	getDefaultLogger().Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	getDefaultLogger().Panicw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	getDefaultLogger().DPanicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	getDefaultLogger().Fatalw(msg, keysAndValues...)
}

func Debugf(msg string, args ...interface{}) {
	getDefaultLogger().Debugf(msg, args...)
}

func Infof(msg string, args ...interface{}) {
	DefaultLogger.Infof(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	getDefaultLogger().Warnf(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	getDefaultLogger().Errorf(msg, args...)
}

func Panicf(msg string, args ...interface{}) {
	DefaultLogger.Panicf(msg, args...)
}

func DPanicf(msg string, args ...interface{}) {
	DefaultLogger.DPanicf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	DefaultLogger.Fatalf(msg, args...)
}

func With(fields ...Field) *Component {
	return DefaultLogger.With(fields...)
}

// withCtx 建议使用
func WithCtx(c context.Context) *Component {
	return DefaultLogger.WithCtx(c)
}
