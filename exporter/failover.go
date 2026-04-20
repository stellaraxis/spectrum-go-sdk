package exporter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	sdklog "go.opentelemetry.io/otel/sdk/log"

	"github.com/stellhub/stellspec-go-sdk/internal/timefmt"
)

// FailoverExporter wraps a primary exporter and persists records locally if export fails.
type FailoverExporter struct {
	primary          sdklog.Exporter
	fallbackFilePath string
	mu               sync.Mutex
}

// NewFailoverExporter creates an exporter that writes failed exports to a local file.
func NewFailoverExporter(primary sdklog.Exporter, fallbackFilePath string) *FailoverExporter {
	return &FailoverExporter{
		primary:          primary,
		fallbackFilePath: fallbackFilePath,
	}
}

// Export sends records through the primary exporter and falls back to a local file on failure.
func (e *FailoverExporter) Export(ctx context.Context, records []sdklog.Record) error {
	err := e.primary.Export(ctx, records)
	if err == nil || len(records) == 0 {
		return err
	}

	// OTLP 重试窗口耗尽后，再把同一批日志顺序追加到本地文件，方便后续补偿和排查。
	if fallbackErr := e.appendFallback(records, err); fallbackErr != nil {
		return errors.Join(err, fallbackErr)
	}
	return err
}

// Shutdown closes the primary exporter.
func (e *FailoverExporter) Shutdown(ctx context.Context) error {
	return e.primary.Shutdown(ctx)
}

// ForceFlush flushes the primary exporter.
func (e *FailoverExporter) ForceFlush(ctx context.Context) error {
	return e.primary.ForceFlush(ctx)
}

func (e *FailoverExporter) appendFallback(records []sdklog.Record, exportErr error) error {
	if e.fallbackFilePath == "" {
		return fmt.Errorf("fallback file path is empty")
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	dir := filepath.Dir(e.fallbackFilePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create fallback directory: %w", err)
		}
	}

	file, err := os.OpenFile(e.fallbackFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open fallback file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	fallbackAt := timefmt.Format(time.Now())
	for _, record := range records {
		payload := map[string]any{
			"fallback_time": fallbackAt,
			"export_error":  exportErr.Error(),
			"record":        buildJSONRecord(record),
		}
		if err := encoder.Encode(payload); err != nil {
			return fmt.Errorf("write fallback record: %w", err)
		}
	}

	return nil
}
