package main

import (
	"context"
	"log"
	"time"

	"github.com/stellaraxis/spectrum-go-sdk/bridge/zapbridge"
	"github.com/stellaraxis/spectrum-go-sdk/config"
	"github.com/stellaraxis/spectrum-go-sdk/sdk"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Default()
	if err := cfg.ApplyEnv(); err != nil {
		log.Fatalf("load config from env failed: %v", err)
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = "spectrum-zap-example"
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

	logger, err := zapbridge.NewLogger(runtime, zapbridge.Options{
		Name:          "examples.zap",
		AddCaller:     true,
		AddStacktrace: true,
		Fields: []zap.Field{
			zap.String("example", "zap"),
			zap.String("component", "demo"),
		},
	})
	if err != nil {
		log.Fatalf("create zap logger failed: %v", err)
	}

	logger.Info("zap example started",
		zap.String("request_id", "req-zap-0001"),
		zap.String("tenant_id", "tenant-a"),
	)

	logger.Warn("slow downstream detected",
		zap.String("dependency", "inventory-service"),
		zap.Duration("cost", 850*time.Millisecond),
	)

	logger.Error("query order failed",
		zap.String("order_id", "order-1001"),
		zap.String("error_code", "ORDER_TIMEOUT"),
		zap.Duration("cost", 1500*time.Millisecond),
	)

	if err := logger.Sync(); err != nil {
		log.Printf("sync logger failed: %v", err)
	}
}
