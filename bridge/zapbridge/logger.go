package zapbridge

import (
	"context"
	"fmt"

	"github.com/stellaraxis/stellspec-go-sdk/requestctx"
	"github.com/stellaraxis/stellspec-go-sdk/sdk"
	otellog "go.opentelemetry.io/otel/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options customizes the zap bridge behavior.
type Options struct {
	Name          string
	Version       string
	SchemaURL     string
	Level         zapcore.LevelEnabler
	AddCaller     bool
	AddStacktrace bool
	Fields        []zap.Field
}

// NewLogger creates a zap logger backed by OpenTelemetry logs.
func NewLogger(runtime *sdk.Runtime, opts Options) (*zap.Logger, error) {
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

	core, err := NewCore(runtime, scopeName, opts.Level, loggerOptions...)
	if err != nil {
		return nil, err
	}
	if len(opts.Fields) > 0 {
		core = core.With(opts.Fields).(*Core)
	}

	zapOptions := make([]zap.Option, 0, 2)
	if runtime.Config().EnableCaller || opts.AddCaller {
		zapOptions = append(zapOptions, zap.AddCaller())
	}
	if runtime.Config().EnableStacktrace || opts.AddStacktrace {
		zapOptions = append(zapOptions, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return zap.New(core, zapOptions...), nil
}

// WithContext binds normalized request context fields to the logger.
func WithContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	if logger == nil {
		return nil
	}

	fields := requestctx.Fields(ctx)
	if len(fields) == 0 {
		return logger
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.String(key, value))
	}
	return logger.With(zapFields...)
}
