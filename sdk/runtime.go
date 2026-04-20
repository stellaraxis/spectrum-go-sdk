package sdk

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/stellaraxis/stellspec-go-sdk/config"
	"github.com/stellaraxis/stellspec-go-sdk/exporter"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Option customizes runtime initialization.
type Option func(*options)

type options struct {
	stdout io.Writer
	stderr io.Writer
}

// WithWriters overrides the development writers used by console output.
func WithWriters(stdout io.Writer, stderr io.Writer) Option {
	return func(o *options) {
		o.stdout = stdout
		o.stderr = stderr
	}
}

// Runtime owns the provider, exporter, and resource lifecycle.
type Runtime struct {
	cfg      config.Config
	provider *sdklog.LoggerProvider
	resource *resource.Resource

	shutdownOnce sync.Once
	shutdownErr  error
}

// New creates a runtime from config and initializes the OTel logger provider.
func New(ctx context.Context, cfg config.Config, optFns ...Option) (*Runtime, error) {
	normalized, err := cfg.Normalize()
	if err != nil {
		return nil, err
	}

	opts := options{}
	for _, optFn := range optFns {
		optFn(&opts)
	}

	res, err := buildResource(normalized)
	if err != nil {
		return nil, err
	}

	exp, err := buildExporter(ctx, normalized, opts)
	if err != nil {
		return nil, err
	}

	processor := buildProcessor(normalized, exp)
	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(processor),
	)

	return &Runtime{
		cfg:      normalized,
		provider: provider,
		resource: res,
	}, nil
}

// Config returns the normalized runtime configuration.
func (r *Runtime) Config() config.Config {
	return r.cfg
}

// Resource returns the OpenTelemetry resource attached to the runtime.
func (r *Runtime) Resource() *resource.Resource {
	return r.resource
}

// LoggerProvider returns the underlying OpenTelemetry logger provider.
func (r *Runtime) LoggerProvider() otellog.LoggerProvider {
	return r.provider
}

// Logger returns a named OpenTelemetry logger.
func (r *Runtime) Logger(name string, opts ...otellog.LoggerOption) otellog.Logger {
	return r.provider.Logger(name, opts...)
}

// Flush forces buffered records to be exported.
func (r *Runtime) Flush(ctx context.Context) error {
	return r.provider.ForceFlush(ctx)
}

// Shutdown flushes buffers and releases exporter resources.
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.shutdownOnce.Do(func() {
		r.shutdownErr = r.provider.Shutdown(ctx)
	})
	return r.shutdownErr
}

func buildProcessor(cfg config.Config, exp sdklog.Exporter) sdklog.Processor {
	if cfg.Output == config.OutputOTLP {
		return sdklog.NewBatchProcessor(
			exp,
			sdklog.WithExportInterval(cfg.BatchTimeout),
			sdklog.WithExportTimeout(cfg.ExportTimeout),
			sdklog.WithExportMaxBatchSize(cfg.MaxBatchSize),
			sdklog.WithMaxQueueSize(cfg.MaxQueueSize),
		)
	}

	return sdklog.NewSimpleProcessor(exp)
}

func buildExporter(ctx context.Context, cfg config.Config, opts options) (sdklog.Exporter, error) {
	if cfg.Output == config.OutputOTLP {
		exp, err := exporter.NewOTLPGRPCExporter(ctx, cfg)
		if err != nil {
			return nil, err
		}
		// 当 OTLP 推送最终失败时，先把日志追加到本地兜底文件，避免 log-agent 不可用时直接丢失。
		return exporter.NewFailoverExporter(exp, cfg.FallbackFilePath), nil
	}

	return exporter.NewConsoleExporter(cfg.Format, cfg.Output, opts.stdout, opts.stderr), nil
}

func buildResource(cfg config.Config) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		attribute.String("service.name", cfg.ServiceName),
	}

	if cfg.ServiceNamespace != "" {
		attrs = append(attrs, attribute.String("service.namespace", cfg.ServiceNamespace))
	}
	if cfg.ServiceVersion != "" {
		attrs = append(attrs, attribute.String("service.version", cfg.ServiceVersion))
	}
	if cfg.ServiceInstanceID != "" {
		attrs = append(attrs, attribute.String("service.instance.id", cfg.ServiceInstanceID))
	}
	if env := strings.TrimSpace(cfg.Environment); env != "" {
		attrs = append(attrs, attribute.String("deployment.environment.name", env))
	}
	if cluster := strings.TrimSpace(cfg.Cluster); cluster != "" {
		attrs = append(attrs, attribute.String("stellar.cluster", cluster))
	}
	if region := strings.TrimSpace(cfg.Region); region != "" {
		attrs = append(attrs, attribute.String("cloud.region", region))
	}
	if zone := strings.TrimSpace(cfg.Zone); zone != "" {
		attrs = append(attrs, attribute.String("cloud.availability_zone", zone))
	}
	if idc := strings.TrimSpace(cfg.IDC); idc != "" {
		attrs = append(attrs, attribute.String("stellar.idc", idc))
	}
	if hostName := strings.TrimSpace(cfg.HostName); hostName != "" {
		attrs = append(attrs, attribute.String("host.name", hostName))
	}
	if hostIP := strings.TrimSpace(cfg.HostIP); hostIP != "" {
		attrs = append(attrs, attribute.String("host.ip", hostIP))
	}
	if nodeName := strings.TrimSpace(cfg.NodeName); nodeName != "" {
		attrs = append(attrs, attribute.String("k8s.node.name", nodeName))
	}
	if namespace := strings.TrimSpace(cfg.K8sNamespace); namespace != "" {
		attrs = append(attrs, attribute.String("k8s.namespace.name", namespace))
	}
	if podName := strings.TrimSpace(cfg.PodName); podName != "" {
		attrs = append(attrs, attribute.String("k8s.pod.name", podName))
	}
	if podIP := strings.TrimSpace(cfg.PodIP); podIP != "" {
		attrs = append(attrs, attribute.String("k8s.pod.ip", podIP))
	}
	if containerName := strings.TrimSpace(cfg.ContainerName); containerName != "" {
		attrs = append(attrs, attribute.String("container.name", containerName))
	}
	for key, value := range cfg.ResourceAttributes {
		if strings.TrimSpace(key) == "" {
			continue
		}
		attrs = append(attrs, attribute.String(key, value))
	}

	base := resource.Default()
	custom := resource.NewSchemaless(attrs...)
	return resource.Merge(base, custom)
}
