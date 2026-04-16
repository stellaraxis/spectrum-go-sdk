package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// ConsoleExporter writes structured records to local writers for development use.
type ConsoleExporter struct {
	mu     sync.Mutex
	stdout io.Writer
	stderr io.Writer
	format string
	output string
	closed bool
}

// NewConsoleExporter creates an exporter that writes records to stdout/stderr.
func NewConsoleExporter(format string, output string, stdout io.Writer, stderr io.Writer) *ConsoleExporter {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	return &ConsoleExporter{
		stdout: stdout,
		stderr: stderr,
		format: format,
		output: output,
	}
}

// Export writes all records as newline-delimited log lines.
func (e *ConsoleExporter) Export(ctx context.Context, records []sdklog.Record) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil
	}

	for _, record := range records {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := e.formatRecord(record)
		if err != nil {
			return err
		}

		if _, err := io.WriteString(e.writerFor(record), line+"\n"); err != nil {
			return err
		}
	}

	return nil
}

// Shutdown marks the exporter as closed.
func (e *ConsoleExporter) Shutdown(context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.closed = true
	return nil
}

// ForceFlush is a no-op because writes happen immediately.
func (e *ConsoleExporter) ForceFlush(context.Context) error {
	return nil
}

func (e *ConsoleExporter) writerFor(record sdklog.Record) io.Writer {
	switch e.output {
	case "stderr":
		return e.stderr
	case "stdout":
		return e.stdout
	default:
		if record.Severity() >= otellog.SeverityError {
			return e.stderr
		}
		return e.stdout
	}
}

func (e *ConsoleExporter) formatRecord(record sdklog.Record) (string, error) {
	if strings.EqualFold(e.format, "console") {
		return formatConsoleRecord(record), nil
	}

	body, err := json.Marshal(buildJSONRecord(record))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func formatConsoleRecord(record sdklog.Record) string {
	parts := []string{
		record.Timestamp().Format(time.RFC3339Nano),
		strings.ToUpper(record.SeverityText()),
		record.Body().String(),
	}

	attrs := map[string]any{}
	record.WalkAttributes(func(kv otellog.KeyValue) bool {
		attrs[kv.Key] = logValueToAny(kv.Value)
		return true
	})
	if len(attrs) > 0 {
		parts = append(parts, fmt.Sprintf("%v", attrs))
	}

	return strings.TrimSpace(strings.Join(parts, " "))
}

func buildJSONRecord(record sdklog.Record) map[string]any {
	payload := map[string]any{
		"timestamp":       record.Timestamp().Format(time.RFC3339Nano),
		"observed_time":   record.ObservedTimestamp().Format(time.RFC3339Nano),
		"severity_number": int(record.Severity()),
		"severity_text":   record.SeverityText(),
		"body":            logValueToAny(record.Body()),
		"attributes":      logAttributesToMap(record),
		"resource":        resourceToMap(record.Resource()),
		"instrumentation": instrumentationToMap(record),
	}

	if traceID := record.TraceID().String(); traceID != "" && traceID != "00000000000000000000000000000000" {
		payload["trace_id"] = traceID
	}
	if spanID := record.SpanID().String(); spanID != "" && spanID != "0000000000000000" {
		payload["span_id"] = spanID
	}
	if eventName := record.EventName(); eventName != "" {
		payload["event_name"] = eventName
	}

	return payload
}

func logAttributesToMap(record sdklog.Record) map[string]any {
	attrs := make(map[string]any, record.AttributesLen())
	record.WalkAttributes(func(kv otellog.KeyValue) bool {
		attrs[kv.Key] = logValueToAny(kv.Value)
		return true
	})
	return attrs
}

func resourceToMap(res interface{ Attributes() []attribute.KeyValue }) map[string]any {
	if res == nil {
		return map[string]any{}
	}

	attrs := res.Attributes()
	result := make(map[string]any, len(attrs))
	for _, attr := range attrs {
		result[string(attr.Key)] = attr.Value.AsInterface()
	}
	return result
}

func instrumentationToMap(record sdklog.Record) map[string]any {
	scope := record.InstrumentationScope()
	result := map[string]any{
		"name": scope.Name,
	}
	if scope.Version != "" {
		result["version"] = scope.Version
	}
	if scope.SchemaURL != "" {
		result["schema_url"] = scope.SchemaURL
	}
	if attrs := scope.Attributes.ToSlice(); len(attrs) > 0 {
		mapped := make(map[string]any, len(attrs))
		for _, attr := range attrs {
			mapped[string(attr.Key)] = attr.Value.AsInterface()
		}
		result["attributes"] = mapped
	}
	return result
}

func logValueToAny(value otellog.Value) any {
	switch value.Kind() {
	case otellog.KindBool:
		return value.AsBool()
	case otellog.KindFloat64:
		return value.AsFloat64()
	case otellog.KindInt64:
		return value.AsInt64()
	case otellog.KindString:
		return value.AsString()
	case otellog.KindBytes:
		return string(value.AsBytes())
	case otellog.KindSlice:
		items := value.AsSlice()
		result := make([]any, 0, len(items))
		for _, item := range items {
			result = append(result, logValueToAny(item))
		}
		return result
	case otellog.KindMap:
		items := value.AsMap()
		result := make(map[string]any, len(items))
		for _, item := range items {
			result[item.Key] = logValueToAny(item.Value)
		}
		return result
	default:
		return value.String()
	}
}
