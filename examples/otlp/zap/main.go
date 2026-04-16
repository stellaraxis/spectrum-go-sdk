package main

import (
	"context"
	"log"
	"time"

	"github.com/stellaraxis/spectrum-go-sdk/bridge/zapbridge"
	"github.com/stellaraxis/spectrum-go-sdk/config"
	"github.com/stellaraxis/spectrum-go-sdk/requestctx"
	"github.com/stellaraxis/spectrum-go-sdk/sdk"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Default()
	cfg.ServiceName = "spectrum-zap-otlp-example"
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

	logger, err := zapbridge.NewLogger(runtime, zapbridge.Options{
		Name:          "examples.otlp.zap",
		AddCaller:     true,
		AddStacktrace: true,
	})
	if err != nil {
		log.Fatalf("create zap logger failed: %v", err)
	}

	ctx := requestctx.WithValues(context.Background(), requestctx.Values{
		RequestID:   "req-prod-zap-0001",
		TenantID:    "tenant-prod-a",
		TraceParent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	})
	logger = zapbridge.WithContext(ctx, logger)

	logger.Info("otlp zap example started",
		zap.String("scene", "prod-otlp"),
	)

	logger.Error("call downstream failed",
		zap.String("dependency", "order-service"),
		zap.String("error_code", "DOWNSTREAM_TIMEOUT"),
	)

	if err := logger.Sync(); err != nil {
		log.Printf("sync logger failed: %v", err)
	}
}
