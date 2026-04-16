package slogbridge

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/stellaraxis/spectrum-go-sdk/internal/otelutil"
	"github.com/stellaraxis/spectrum-go-sdk/internal/severity"
	"github.com/stellaraxis/spectrum-go-sdk/requestctx"
	"github.com/stellaraxis/spectrum-go-sdk/sdk"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
)

// Options customizes the slog bridge behavior.
type Options struct {
	Name      string
	Version   string
	SchemaURL string
	Level     slog.Leveler
	AddSource bool
	Attrs     []slog.Attr
}

// Handler bridges slog records into OpenTelemetry logs.
type Handler struct {
	runtime   *sdk.Runtime
	logger    otellog.Logger
	level     slog.Leveler
	addSource bool
	attrs     []slog.Attr
	groups    []string
}

// NewHandler creates a slog handler backed by the runtime.
func NewHandler(runtime *sdk.Runtime, opts Options) (*Handler, error) {
	if runtime == nil {
		return nil, fmt.Errorf("runtime is required")
	}

	scopeName := opts.Name
	if scopeName == "" {
		scopeName = runtime.Config().ServiceName
	}

	loggerOptions := make([]otellog.LoggerOption, 0, 2)
	if opts.Version != "" {
		loggerOptions = append(loggerOptions, otellog.WithInstrumentationVersion(opts.Version))
	}
	if opts.SchemaURL != "" {
		loggerOptions = append(loggerOptions, otellog.WithSchemaURL(opts.SchemaURL))
	}

	level := opts.Level
	if level == nil {
		level = parseLevel(runtime.Config().Level)
	}

	handler := &Handler{
		runtime:   runtime,
		logger:    runtime.Logger(scopeName, loggerOptions...),
		level:     level,
		addSource: runtime.Config().EnableCaller || opts.AddSource,
		attrs:     append([]slog.Attr(nil), opts.Attrs...),
	}

	return handler, nil
}

// Enabled reports whether the level should be emitted.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle emits the record as an OpenTelemetry log record.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	otelRecord := otellog.Record{}
	if !record.Time.IsZero() {
		otelRecord.SetTimestamp(record.Time)
		otelRecord.SetObservedTimestamp(record.Time)
	}
	otelRecord.SetSeverity(severity.FromSlog(record.Level))
	otelRecord.SetSeverityText(severity.TextFromSlog(record.Level))
	otelRecord.SetBody(otellog.StringValue(record.Message))
	setTraceContext(ctx, &otelRecord)

	attrs := make([]otellog.KeyValue, 0, len(h.attrs)+record.NumAttrs()+4)
	attrs = append(attrs, attrsToOTel(h.groups, h.attrs...)...)
	attrs = append(attrs, contextAttrs(ctx)...)
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attrToOTel(h.groups, attr)...)
		return true
	})

	if h.addSource && record.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{record.PC})
		frame, _ := frames.Next()
		if frame.File != "" {
			attrs = append(attrs,
				otellog.String("code.filepath", frame.File),
				otellog.Int("code.lineno", frame.Line),
			)
		}
		if frame.Function != "" {
			attrs = append(attrs, otellog.String("code.function", frame.Function))
		}
	}

	otelRecord.AddAttributes(attrs...)
	h.logger.Emit(ctx, otelRecord)
	return nil
}

// WithAttrs returns a cloned handler with additional attrs.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	cloned := h.clone()
	cloned.attrs = append(cloned.attrs, attrs...)
	return cloned
}

// WithGroup returns a cloned handler with one more attribute group.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	cloned := h.clone()
	cloned.groups = append(cloned.groups, name)
	return cloned
}

func (h *Handler) clone() *Handler {
	return &Handler{
		runtime:   h.runtime,
		logger:    h.logger,
		level:     h.level,
		addSource: h.addSource,
		attrs:     append([]slog.Attr(nil), h.attrs...),
		groups:    append([]string(nil), h.groups...),
	}
}

func parseLevel(level string) slog.Leveler {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func setTraceContext(ctx context.Context, record *otellog.Record) {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return
	}

	record.AddAttributes(
		otellog.String("trace.id", spanContext.TraceID().String()),
		otellog.String("span.id", spanContext.SpanID().String()),
		otellog.String("trace.flags", spanContext.TraceFlags().String()),
	)
}

func attrsToOTel(groups []string, attrs ...slog.Attr) []otellog.KeyValue {
	result := make([]otellog.KeyValue, 0, len(attrs))
	for _, attr := range attrs {
		result = append(result, attrToOTel(groups, attr)...)
	}
	return result
}

func attrToOTel(groups []string, attr slog.Attr) []otellog.KeyValue {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return nil
	}

	if attr.Value.Kind() == slog.KindGroup {
		groupAttrs := attr.Value.Group()
		if attr.Key == "" {
			return attrsToOTel(groups, groupAttrs...)
		}
		nested := attrsToOTel(append(groups, attr.Key), groupAttrs...)
		if len(nested) == 0 {
			return nil
		}
		return []otellog.KeyValue{
			{
				Key:   joinGroups(groups, attr.Key),
				Value: otellog.MapValue(nested...),
			},
		}
	}

	return []otellog.KeyValue{
		{
			Key:   joinGroups(groups, attr.Key),
			Value: otelutil.AnyToValue(slogValueToAny(attr.Value)),
		},
	}
}

func joinGroups(groups []string, key string) string {
	if len(groups) == 0 {
		return key
	}
	qualified := append(append([]string(nil), groups...), key)
	return strings.Join(qualified, ".")
}

func slogValueToAny(value slog.Value) any {
	switch value.Kind() {
	case slog.KindBool:
		return value.Bool()
	case slog.KindDuration:
		return value.Duration()
	case slog.KindFloat64:
		return value.Float64()
	case slog.KindInt64:
		return value.Int64()
	case slog.KindString:
		return value.String()
	case slog.KindTime:
		return value.Time()
	case slog.KindUint64:
		return value.Uint64()
	case slog.KindAny:
		return value.Any()
	default:
		return value.String()
	}
}

func contextAttrs(ctx context.Context) []otellog.KeyValue {
	fields := requestctx.Fields(ctx)
	if len(fields) == 0 {
		return nil
	}

	attrs := make([]otellog.KeyValue, 0, len(fields))
	for key, value := range fields {
		attrs = append(attrs, otellog.String(key, value))
	}
	return attrs
}
