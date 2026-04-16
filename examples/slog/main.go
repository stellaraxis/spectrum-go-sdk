package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/stellaraxis/spectrum-go-sdk/bridge/slogbridge"
	"github.com/stellaraxis/spectrum-go-sdk/config"
	"github.com/stellaraxis/spectrum-go-sdk/sdk"
)

func main() {
	cfg := config.Default()
	if err := cfg.ApplyEnv(); err != nil {
		log.Fatalf("load config from env failed: %v", err)
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = "spectrum-slog-example"
	}
	if cfg.ServiceNamespace == "" {
		cfg.ServiceNamespace = "stellar.examples"
	}
	if cfg.ServiceVersion == "" {
		cfg.ServiceVersion = "1.0.0"
	}
	if cfg.Environment == "" {
		cfg.Environment = "dev"
	}
	if cfg.Output == "" || cfg.Output == config.OutputOTLP {
		cfg.Output = config.OutputConsole
	}
	if cfg.Format == "" {
		cfg.Format = config.FormatConsole
	}
	cfg.Development = true

	cfg, err := cfg.Normalize()
	if err != nil {
		log.Fatalf("normalize config failed: %v", err)
	}

	runtime, err := sdk.New(context.Background(), cfg)
	if err != nil {
		log.Fatalf("init runtime failed: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := runtime.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown runtime failed: %v", err)
		}
	}()

	handler, err := slogbridge.NewHandler(runtime, slogbridge.Options{
		Name:      "examples.slog",
		AddSource: true,
		Attrs: []slog.Attr{
			slog.String("example", "slog"),
			slog.String("component", "demo"),
		},
	})
	if err != nil {
		log.Fatalf("create slog handler failed: %v", err)
	}

	logger := slog.New(handler).With(
		slog.String("tenant_id", "tenant-b"),
		slog.String("request_id", "req-slog-0001"),
	)

	logger.Info("slog example started",
		slog.String("scene", "boot"),
	)

	logger.Warn("slow downstream detected",
		slog.String("dependency", "payment-service"),
		slog.Duration("cost", 920*time.Millisecond),
	)

	logger.Error("push invoice failed",
		slog.String("invoice_id", "invoice-2001"),
		slog.String("reason", "timeout"),
		slog.Int64("retry", 2),
	)
}
