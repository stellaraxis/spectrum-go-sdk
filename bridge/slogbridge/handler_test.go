package slogbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"unicode/utf8"

	"github.com/stellhub/stellspec-go-sdk/config"
	"github.com/stellhub/stellspec-go-sdk/internal/logbody"
	"github.com/stellhub/stellspec-go-sdk/requestctx"
	"github.com/stellhub/stellspec-go-sdk/sdk"
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

func TestHandlerTruncatesLongBody(t *testing.T) {
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

	handler, err := NewHandler(runtime, Options{Name: "test-slog-truncate"})
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	message := bytes.Repeat([]byte("星"), logbody.MaxBytes+128)
	slog.New(handler).Info(string(message))

	if err := runtime.Flush(context.Background()); err != nil {
		t.Fatalf("flush runtime: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal log line: %v", err)
	}

	body, ok := payload["body"].(string)
	if !ok {
		t.Fatalf("body should be string: %T", payload["body"])
	}
	if len(body) > logbody.MaxBytes {
		t.Fatalf("body should be truncated to <= %d bytes, got %d", logbody.MaxBytes, len(body))
	}
	if !utf8.ValidString(body) {
		t.Fatal("body should remain valid utf-8 after truncation")
	}

	attrs, ok := payload["attributes"].(map[string]any)
	if !ok {
		t.Fatal("attributes should be a map")
	}
	if attrs["log.body_truncated"] != true {
		t.Fatalf("unexpected truncation flag: %v", attrs["log.body_truncated"])
	}
	if attrs["log.body_original_size"] != float64(len(string(message))) {
		t.Fatalf("unexpected original size: %v", attrs["log.body_original_size"])
	}
	if attrs["log.body_max_size"] != float64(logbody.MaxBytes) {
		t.Fatalf("unexpected max size: %v", attrs["log.body_max_size"])
	}
}
