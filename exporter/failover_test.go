package exporter

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func TestFailoverExporterWritesLocalFileWhenPrimaryExportFails(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "fallback", "failed-log.jsonl")
	wantErr := errors.New("otlp export failed")
	exp := NewFailoverExporter(stubExporter{exportErr: wantErr}, path)

	record := sdklog.Record{}
	record.SetTimestamp(time.Unix(1710000000, 0))
	record.SetObservedTimestamp(time.Unix(1710000000, 0))
	record.SetSeverity(otellog.SeverityInfo)
	record.SetSeverityText("INFO")
	record.SetBody(otellog.StringValue("push order failed"))

	err := exp.Export(context.Background(), []sdklog.Record{record})
	if !errors.Is(err, wantErr) {
		t.Fatalf("unexpected export error: %v", err)
	}

	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read fallback file: %v", readErr)
	}

	text := string(content)
	if !strings.Contains(text, "\"export_error\":\"otlp export failed\"") {
		t.Fatalf("fallback file does not contain export error: %s", text)
	}
	if !strings.Contains(text, "\"body\":\"push order failed\"") {
		t.Fatalf("fallback file does not contain record body: %s", text)
	}
}

func TestFailoverExporterDoesNotWriteLocalFileWhenPrimaryExportSucceeds(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "fallback", "success-log.jsonl")
	exp := NewFailoverExporter(stubExporter{}, path)

	record := sdklog.Record{}
	record.SetBody(otellog.StringValue("ok"))

	if err := exp.Export(context.Background(), []sdklog.Record{record}); err != nil {
		t.Fatalf("export should succeed: %v", err)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("fallback file should not exist, stat err: %v", err)
	}
}

type stubExporter struct {
	exportErr error
}

func (e stubExporter) Export(context.Context, []sdklog.Record) error {
	return e.exportErr
}

func (e stubExporter) Shutdown(context.Context) error {
	return nil
}

func (e stubExporter) ForceFlush(context.Context) error {
	return nil
}
