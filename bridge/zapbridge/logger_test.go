package zapbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stellaraxis/spectrum-go-sdk/config"
	"github.com/stellaraxis/spectrum-go-sdk/requestctx"
	"github.com/stellaraxis/spectrum-go-sdk/sdk"
)

func TestNewLogger(t *testing.T) {
	stdout := new(bytes.Buffer)
	runtime, err := sdk.New(context.Background(), config.Config{
		ServiceName: "user-service",
		Environment: "dev",
		Format:      config.FormatJSON,
	}, sdk.WithWriters(stdout, stdout))
	if err != nil {
		t.Fatalf("new runtime: %v", err)
	}
	t.Cleanup(func() {
		_ = runtime.Shutdown(context.Background())
	})

	logger, err := NewLogger(runtime, Options{Name: "test-zap"})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	ctx := requestctx.WithValues(context.Background(), requestctx.Values{
		RequestID:   "req-1",
		TenantID:    "tenant-a",
		TraceParent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	})
	logger = WithContext(ctx, logger)
	logger.Info("hello")
	if err := logger.Sync(); err != nil {
		t.Fatalf("sync logger: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal log line: %v", err)
	}
	attrs, ok := payload["attributes"].(map[string]any)
	if !ok {
		t.Fatal("attributes should be a map")
	}
	if attrs["request_id"] != "req-1" {
		t.Fatalf("unexpected request id: %v", attrs["request_id"])
	}
	if attrs["tenant_id"] != "tenant-a" {
		t.Fatalf("unexpected tenant id: %v", attrs["tenant_id"])
	}
	if attrs["traceparent"] == "" {
		t.Fatal("traceparent should not be empty")
	}
}
