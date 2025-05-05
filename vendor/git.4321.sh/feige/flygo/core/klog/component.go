package klog

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"git.4321.sh/feige/flygo/core/kcfg"
	"git.4321.sh/feige/flygo/core/utils/xcolor"
)

const (
	DebugLevel = zap.DebugLevel
	InfoLevel  = zap.InfoLevel
	WarnLevel  = zap.WarnLevel
	ErrorLevel = zap.ErrorLevel
	PanicLevel = zap.PanicLevel
	FatalLevel = zap.FatalLevel
)

type (
	Field     = zap.Field
	Level     = zapcore.Level
	Component struct {
		name          string
		desugar       *zap.Logger
		lv            *zap.AtomicLevel
		config        *Config
		sugar         *zap.SugaredLogger
		asyncStopFunc func() error
	}
)

var (
	String     = zap.String
	Any        = zap.Any
	Int64      = zap.Int64
	Int        = zap.Int
	Int32      = zap.Int32
	Uint       = zap.Uint
	Duration   = zap.Duration
	Durationp  = zap.Durationp
	Object     = zap.Object
	Namespace  = zap.Namespace
	Reflect    = zap.Reflect
	Skip       = zap.Skip()
	ByteString = zap.ByteString
)

func newStderrCore(config *Config, lv zap.AtomicLevel) (zapcore.Core, CloseFunc) {
	return zapcore.NewCore(zapcore.NewJSONEncoder(*config.encoderConfig), os.Stderr, lv), noopCloseFunc
}

// 文件输出 log
func newRotateFileCore(config *Config, lv zap.AtomicLevel) (zapcore.Core, CloseFunc) {
	cf := noopCloseFunc
	var ws = zapcore.AddSync(newRotate(config))
	if config.Debug {
		ws = zap.CombineWriteSyncers(os.Stdout, ws)
	}
	if config.EnableAsync {
		ws, cf = bufferWriteSyncer(ws, config.FlushBufferSize, config.FlushBufferInterval)
	}
	core := zapcore.NewCore(
		func() zapcore.Encoder {
			if config.Debug {
				return zapcore.NewConsoleEncoder(*config.encoderConfig)
			}
			return zapcore.NewJSONEncoder(*config.encoderConfig)
		}(),
		ws,
		lv,
	)
	return core, cf
}

func newCore(config *Config, lv zap.AtomicLevel) (zapcore.Core, CloseFunc) {
	switch config.Writer {
	case writerRotateFile:
		return newRotateFileCore(config, lv)
	case writerStderr:
		return newStderrCore(config, lv)
	default:
		panic("unsupported writer")
	}
}

func newLogger(name string, config *Config) *Component {
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if config.EnableAddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}
	if len(config.fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(config.fields...))
	}

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(config.Level)); err != nil {
		panic(err)
	}
	core, asyncStopFunc := newCore(config, lv)
	zapLogger := zap.New(core, zapOptions...)
	return &Component{
		desugar:       zapLogger,
		lv:            &lv,
		config:        config,
		sugar:         zapLogger.Sugar(),
		name:          name,
		asyncStopFunc: asyncStopFunc,
	}
}

func (logger *Component) AutoLevel(confKey string) {
	kcfg.OnChange(func(config *kcfg.Configuration) {
		lvText := strings.ToLower(config.GetString(confKey))
		if lvText != "" {
			logger.Info("update level", String("level", lvText), String("name", logger.config.Name))
			_ = logger.lv.UnmarshalText([]byte(lvText))
		}
	})
}

func (logger *Component) SetLevel(lv Level) {
	logger.lv.SetLevel(lv)
}

func (logger *Component) Flush() error {
	if logger.asyncStopFunc != nil {
		if err := logger.asyncStopFunc(); err != nil {
			return err
		}
	}

	_ = logger.desugar.Sync()
	return nil
}

func defaultZapConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "body",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeUnixMicroEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

//
func defaultDebugConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "body",
		StacktraceKey: "stack",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		//EncodeLevel:    debugEncodeLevel,
		EncodeTime:     timeUnixMicroEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func debugEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = xcolor.Red
	switch lv {
	case zapcore.DebugLevel:
		colorize = xcolor.Blue
	case zapcore.InfoLevel:
		colorize = xcolor.Green
	case zapcore.WarnLevel:
		colorize = xcolor.Yellow
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize = xcolor.Red
	default:
	}
	enc.AppendString(colorize(lv.CapitalString()))
}

// 日志时间格式 - 时间戳(秒)
func timestampEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}

// 日志时间格式 - 时间戳（毫秒）
func timeUnixMicroEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix() * 1000)
}

// 日志时间格式 - YYYY-MM-DD hh:ii:ss
func dateTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func (logger *Component) IsDebugMode() bool {
	return logger.config.Debug
}

func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-32s", msg)
}

func (logger *Component) Debug(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Debug(msg, fields...)
}

func (logger *Component) Debugw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Debugw(msg, keysAndValues...)
}

func sprintf(template string, args ...interface{}) string {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	return msg
}

func (logger *Component) StdLog() *log.Logger {
	return zap.NewStdLog(logger.desugar)
}

func (logger *Component) Debugf(template string, args ...interface{}) {
	logger.sugar.Debugw(sprintf(template, args...))
}

func (logger *Component) Info(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Info(msg, fields...)
}

func (logger *Component) Infow(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Infow(msg, keysAndValues...)
}

func (logger *Component) Infof(template string, args ...interface{}) {
	logger.sugar.Infof(sprintf(template, args...))
}

func (logger *Component) Warn(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Warn(msg, fields...)
}

func (logger *Component) Warnw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Warnw(msg, keysAndValues...)
}

func (logger *Component) Warnf(template string, args ...interface{}) {
	logger.sugar.Warnf(sprintf(template, args...))
}

func (logger *Component) Error(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Error(msg, fields...)
}

func (logger *Component) Errorw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Errorw(msg, keysAndValues...)
}

func (logger *Component) Errorf(template string, args ...interface{}) {
	logger.sugar.Errorf(sprintf(template, args...))
}

func (logger *Component) Panic(msg string, fields ...Field) {
	panicDetail(msg, fields...)
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Panic(msg, fields...)
}

func (logger *Component) Panicw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Panicw(msg, keysAndValues...)
}

func (logger *Component) Panicf(template string, args ...interface{}) {
	logger.sugar.Panicf(sprintf(template, args...))
}

func (logger *Component) DPanic(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
	}
	logger.desugar.DPanic(msg, fields...)
}

func (logger *Component) DPanicw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.DPanicw(msg, keysAndValues...)
}

func (logger *Component) DPanicf(template string, args ...interface{}) {
	logger.sugar.DPanicf(sprintf(template, args...))
}

func (logger *Component) Fatal(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		//msg = normalizeMessage(msg)
		return
	}
	logger.desugar.Fatal(msg, fields...)
}

func (logger *Component) Fatalw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Fatalw(msg, keysAndValues...)
}

func (logger *Component) Fatalf(template string, args ...interface{}) {
	logger.sugar.Fatalf(sprintf(template, args...))
}

func panicDetail(msg string, fields ...Field) {
	enc := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(enc)
	}

	fmt.Printf("%s: \n    %s: %s\n", xcolor.Red("panic"), xcolor.Red("msg"), msg)
	if _, file, line, ok := runtime.Caller(3); ok {
		fmt.Printf("    %s: %s:%d\n", xcolor.Red("loc"), file, line)
	}
	for key, val := range enc.Fields {
		fmt.Printf("    %s: %s\n", xcolor.Red(key), fmt.Sprintf("%+v", val))
	}
}

// With 添加日志字段
func (logger *Component) With(fields ...Field) *Component {
	desugarLogger := logger.desugar.With(fields...)
	return &Component{
		desugar: desugarLogger,
		lv:      logger.lv,
		sugar:   desugarLogger.Sugar(),
		config:  logger.config,
	}
}

// WithCtx ...
func (logger *Component) WithCtx(ctx context.Context) *Component {
	fields := make([]Field, 0)
	fields = append(fields,
		FieldCtxUid(ctx),
	)
	return logger.With(fields...)
}

func (logger *Component) WithCallerSkip(callerSkip int, fields ...Field) *Component {
	logger.config.CallerSkip = callerSkip
	desugarLogger := logger.desugar.WithOptions(zap.AddCallerSkip(callerSkip)).With(fields...)
	return &Component{
		desugar: desugarLogger,
		lv:      logger.lv,
		sugar:   desugarLogger.Sugar(),
		config:  logger.config,
	}
}

func (logger *Component) ConfigDir() string {
	return logger.config.Dir
}

func (logger *Component) ConfigName() string {
	return logger.config.Name
}
