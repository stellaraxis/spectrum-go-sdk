package slogbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stellaraxis/spectrum-go-sdk/config"
	"github.com/stellaraxis/spectrum-go-sdk/requestctx"
	"github.com/stellaraxis/spectrum-go-sdk/sdk"
)

func TestNewHandler(t *testing.T) {
	stdout := new(bytes.Buffer)
	runtime, err := sdk.New(context.Background(), config.Config{
		ServiceName: "billing-service",
		Environment: "dev",
		Format:      config.FormatJSON,
	}, sdk.WithWriters(stdout, stdout))
	if err != nil {
		t.Fatalf("new runtime: %v", err)
	}
	t.Cleanup(func() {
		_ = runtime.Shutdown(context.Background())
	})

	handler, err := NewHandler(runtime, Options{Name: "test-slog"})
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	logger := slog.New(handler)
	ctx := requestctx.WithValues(context.Background(), requestctx.Values{
		RequestID:   "req-1",
		TenantID:    "tenant-a",
		TraceParent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	})
	logger.InfoContext(ctx, "hello")

	if err := runtime.Flush(context.Background()); err != nil {
		t.Fatalf("flush runtime: %v", err)
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
