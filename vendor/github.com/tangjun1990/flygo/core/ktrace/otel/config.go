package otel

import (
	"context"

	"github.com/tangjun1990/flygo/core/kapp"
	"github.com/tangjun1990/flygo/core/kcfg"
	"github.com/tangjun1990/flygo/core/klog"

	jaegerv2 "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// Config ...
type Config struct {
	ServiceName  string
	OtelType     string  // type: otlp ,jaeger
	Fraction     float64 // 采样率： 默认0不会采集
	PanicOnError bool
	options      []tracesdk.TracerProviderOption
	Jaeger       jaegerConfig // otel jaeger 配置
	Otlp         otlpConfig   // otel otlp 配置
}

// otlpConfig otlp上报协议配置
type otlpConfig struct {
	Endpoint   string                 // oltp endpoint
	Headers    map[string]string      // 默认提供一个 请求头的参数配置
	options    []otlptracegrpc.Option // 预留自定义配置：   例如 grpc WithGRPCConn
	resOptions []resource.Option      // res 预留自定以配置
}

// jaegerConfig jaeger上报协议配置
type jaegerConfig struct {
	EndpointType      string // type: agent,collector
	AgentHost         string // agent host
	AgentPort         string // agent port
	CollectorEndpoint string // collector endpoint
	CollectorUser     string // collector user
	CollectorPassword string // collector password
}

// Load parses container configuration from configuration provider, such as a toml file,
// then use the configuration to construct a component container.
func Load(key string) *Config {
	var config = DefaultConfig()
	if err := kcfg.UnmarshalKey(key, config); err != nil {
		panic(err)
	}
	return config
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		ServiceName: kapp.Name(),
		Jaeger: jaegerConfig{
			AgentHost:         "localhost",
			AgentPort:         "6831",
			CollectorEndpoint: "http://localhost:14268/api/traces",
			CollectorUser:     "",
			CollectorPassword: "",
			EndpointType:      "collector",
		},
		Otlp: otlpConfig{
			Endpoint: "localhost:4317",
		},
		OtelType:     "otlp",
		PanicOnError: true,
	}
}

// WithTracerProviderOption ...
func (config *Config) WithTracerProviderOption(options ...tracesdk.TracerProviderOption) *Config {
	config.options = append(config.options, options...)
	return config
}

// WithOtlpTraceGrpcOption 自定义otlp Option
func (config *Config) WithOtlpTraceGrpcOption(options ...otlptracegrpc.Option) *Config {
	config.Otlp.options = append(config.Otlp.options, options...)
	return config
}

// WithOtlpResourceOption 自定义otlp resource Option
func (config *Config) WithOtlpResourceOption(options ...resource.Option) *Config {
	config.Otlp.resOptions = append(config.Otlp.resOptions, options...)
	return config
}

// Build ...
func (config *Config) Build(options ...Option) trace.TracerProvider {
	for _, option := range options {
		option(config)
	}
	switch config.OtelType {
	case "otlp":
		return config.buildOtlpTP()
	case "jaeger":
		return config.buildJaegerTP()
	default:
		klog.Panic("otel type error", klog.FieldName(config.OtelType))
	}
	return nil
}

func (config *Config) buildJaegerTP() trace.TracerProvider {
	var endpoint jaegerv2.EndpointOption
	switch config.Jaeger.EndpointType {
	case "agent":
		// Create the Jaeger exporter
		endpoint = jaegerv2.WithAgentEndpoint(
			jaegerv2.WithAgentHost(config.Jaeger.AgentHost),
			jaegerv2.WithAgentPort(config.Jaeger.AgentPort),
		)
	case "collector":
		endpoint = jaegerv2.WithCollectorEndpoint(
			jaegerv2.WithEndpoint(config.Jaeger.CollectorEndpoint),
			jaegerv2.WithUsername(config.Jaeger.CollectorUser),
			jaegerv2.WithPassword(config.Jaeger.CollectorPassword),
		)
	default:
		klog.Panic("jaeger type error", klog.FieldName(config.Jaeger.EndpointType))
	}

	exp, err := jaegerv2.New(endpoint)
	if err != nil {
		return nil
	}
	options := []tracesdk.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(config.Fraction))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(config.ServiceName),
		)),
	}
	options = append(options, config.options...)
	tp := tracesdk.NewTracerProvider(options...)
	return tp
}

func (config *Config) buildOtlpTP() trace.TracerProvider {
	// otlp exporter
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),                     // WithInsecure disables client transport security for the exporter's gRPC
		otlptracegrpc.WithHeaders(config.Otlp.Headers),   // WithHeaders will send the provided headers with each gRPC requests.
		otlptracegrpc.WithEndpoint(config.Otlp.Endpoint), // WithEndpoint sets the target endpoint the exporter will connect to. If unset, localhost:4317 will be used as a default.
		// otlptracegrpc.WithDialOption(grpc.WithBlock()), //默认不设置 同步状态，会产生阻塞等待 Ready
	}
	options = append(options, config.Otlp.options...)
	traceClient := otlptracegrpc.NewClient(options...)
	ctx := context.Background()
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		klog.Error("otlp exporter error", klog.FieldErr(err))
		return nil
	}

	// res
	resOptions := []resource.Option{
		resource.WithTelemetrySDK(), // WithTelemetrySDK adds TelemetrySDK version info to the configured resource.
		resource.WithHost(),         // WithHost adds attributes from the host to the configured resource.
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(config.ServiceName),
		),
	}
	resOptions = append(resOptions, config.Otlp.resOptions...)
	res, err := resource.New(ctx, resOptions...)
	if err != nil {
		klog.Error("otlp resource New error", klog.FieldErr(err))
		return nil
	}

	// tp
	tpOptions := []tracesdk.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(config.Fraction))),
		// WithSpanProcessor registers the SpanProcessor with a TracerProvider.
		tracesdk.WithSpanProcessor(tracesdk.NewBatchSpanProcessor(traceExp)),
		// Record information about this application in an Resource.
		tracesdk.WithResource(res),
	}
	tpOptions = append(tpOptions, config.options...)
	tp := tracesdk.NewTracerProvider(tpOptions...)
	return tp
}

// Stop 停止
func (config *Config) Stop() error {
	return nil
}
