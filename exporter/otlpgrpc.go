package exporter

import (
	"context"

	"github.com/stellhub/stellspec-go-sdk/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewOTLPGRPCExporter creates a production exporter backed by OTLP/gRPC.
func NewOTLPGRPCExporter(ctx context.Context, cfg config.Config) (sdklog.Exporter, error) {
	options := []otlploggrpc.Option{
		// 这里会把日志通过 OTLP/gRPC 发送到 cfg.Endpoint，生产环境通常就是本机 log-agent。
		otlploggrpc.WithEndpoint(cfg.Endpoint),
		otlploggrpc.WithTimeout(cfg.ExportTimeout),
	}
	// exporter 级别重试用于处理 log-agent 短暂不可用、连接抖动等瞬时错误；
	// 如果重试窗口耗尽，则交由外层失败落盘逻辑继续兜底。
	options = append(options, otlploggrpc.WithRetry(otlploggrpc.RetryConfig{
		Enabled:         *cfg.Retry.Enabled,
		InitialInterval: *cfg.Retry.InitialInterval,
		MaxInterval:     *cfg.Retry.MaxInterval,
		MaxElapsedTime:  *cfg.Retry.MaxElapsedTime,
	}))

	if cfg.Insecure {
		options = append(options, otlploggrpc.WithInsecure())
	}
	if len(cfg.Headers) > 0 {
		options = append(options, otlploggrpc.WithHeaders(cfg.Headers))
	}

	return otlploggrpc.New(ctx, options...)
}
