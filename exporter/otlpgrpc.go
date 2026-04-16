package exporter

import (
	"context"

	"github.com/stellaraxis/spectrum-go-sdk/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewOTLPGRPCExporter creates a production exporter backed by OTLP/gRPC.
func NewOTLPGRPCExporter(ctx context.Context, cfg config.Config) (sdklog.Exporter, error) {
	options := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(cfg.Endpoint),
		otlploggrpc.WithTimeout(cfg.ExportTimeout),
	}

	if cfg.Insecure {
		options = append(options, otlploggrpc.WithInsecure())
	}
	if len(cfg.Headers) > 0 {
		options = append(options, otlploggrpc.WithHeaders(cfg.Headers))
	}

	return otlploggrpc.New(ctx, options...)
}
