package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/stellaraxis/stellspec-go-sdk/bridge/slogbridge"
	"github.com/stellaraxis/stellspec-go-sdk/config"
	"github.com/stellaraxis/stellspec-go-sdk/requestctx"
	"github.com/stellaraxis/stellspec-go-sdk/sdk"
)

func main() {
	cfg := config.Default()
	cfg.ServiceName = "stellspec-slog-otlp-example"
	cfg.ServiceNamespace = "stellar.examples"
	cfg.ServiceVersion = "1.0.0"
	cfg.Environment = "prod"
	cfg.Development = false
	cfg.Output = config.OutputOTLP
	cfg.Format = config.FormatJSON
	cfg.Endpoint = "localhost:4317"
	cfg.Insecure = true

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
		Name:      "examples.otlp.slog",
		AddSource: true,
	})
	if err != nil {
		log.Fatalf("create slog handler failed: %v", err)
	}

	ctx := requestctx.WithValues(context.Background(), requestctx.Values{
		RequestID:   "req-prod-slog-0001",
		TenantID:    "tenant-prod-b",
		TraceParent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	})
	logger := slog.New(handler)

	logger.InfoContext(ctx, "otlp slog example started",
		slog.String("scene", "prod-otlp"),
	)

	logger.ErrorContext(ctx, "push invoice failed",
		slog.String("dependency", "invoice-service"),
		slog.String("error_code", "PUSH_TIMEOUT"),
	)
}
