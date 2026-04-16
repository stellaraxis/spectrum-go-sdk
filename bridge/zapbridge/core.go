package zapbridge

import (
	"context"

	"github.com/stellaraxis/spectrum-go-sdk/internal/otelutil"
	"github.com/stellaraxis/spectrum-go-sdk/internal/severity"
	"github.com/stellaraxis/spectrum-go-sdk/sdk"
	otellog "go.opentelemetry.io/otel/log"
	"go.uber.org/zap/zapcore"
)

// Core is a zapcore.Core implementation that emits OTel log records.
type Core struct {
	runtime *sdk.Runtime
	logger  otellog.Logger
	level   zapcore.LevelEnabler
	fields  []zapcore.Field
}

// NewCore creates a zap core backed by the runtime's logger provider.
func NewCore(runtime *sdk.Runtime, name string, level zapcore.LevelEnabler, opts ...otellog.LoggerOption) (*Core, error) {
	if level == nil {
		parsed, err := zapcore.ParseLevel(runtime.Config().Level)
		if err != nil {
			return nil, err
		}
		level = parsed
	}

	return &Core{
		runtime: runtime,
		logger:  runtime.Logger(name, opts...),
		level:   level,
	}, nil
}

// Enabled reports whether the core should emit the given level.
func (c *Core) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

// With returns a new core with additional default fields.
func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	cloned := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	cloned = append(cloned, c.fields...)
	cloned = append(cloned, fields...)

	return &Core{
		runtime: c.runtime,
		logger:  c.logger,
		level:   c.level,
		fields:  cloned,
	}
}

// Check adds the core to the checked entry when the level is enabled.
func (c *Core) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

// Write converts the zap entry into an OTel log record and emits it.
func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	record := otellog.Record{}
	record.SetTimestamp(entry.Time)
	record.SetObservedTimestamp(entry.Time)
	record.SetSeverity(severity.FromZap(entry.Level))
	record.SetSeverityText(severity.TextFromZap(entry.Level))
	record.SetBody(otellog.StringValue(entry.Message))

	attrs, err := fieldsToAttributes(c.fields, fields)
	if err != nil {
		return err
	}
	if entry.LoggerName != "" {
		attrs = append(attrs, otellog.String("logger.name", entry.LoggerName))
	}
	if entry.Caller.Defined {
		attrs = append(attrs,
			otellog.String("code.filepath", entry.Caller.File),
			otellog.Int("code.lineno", entry.Caller.Line),
		)
		if entry.Caller.Function != "" {
			attrs = append(attrs, otellog.String("code.function", entry.Caller.Function))
		}
	}
	if entry.Stack != "" {
		attrs = append(attrs, otellog.String("exception.stacktrace", entry.Stack))
	}
	record.AddAttributes(attrs...)

	c.logger.Emit(context.Background(), record)
	return nil
}

// Sync flushes the runtime buffers.
func (c *Core) Sync() error {
	return c.runtime.Flush(context.Background())
}

func fieldsToAttributes(base []zapcore.Field, extra []zapcore.Field) ([]otellog.KeyValue, error) {
	allFields := make([]zapcore.Field, 0, len(base)+len(extra))
	allFields = append(allFields, base...)
	allFields = append(allFields, extra...)

	encoder := zapcore.NewMapObjectEncoder()
	for _, field := range allFields {
		field.AddTo(encoder)
	}

	return otelutil.MapToAttributes(encoder.Fields), nil
}
