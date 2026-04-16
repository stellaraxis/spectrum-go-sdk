package severity

import (
	"log/slog"
	"strings"

	otellog "go.opentelemetry.io/otel/log"
	"go.uber.org/zap/zapcore"
)

// FromZap converts a zap level to OpenTelemetry severity.
func FromZap(level zapcore.Level) otellog.Severity {
	switch {
	case level <= zapcore.DebugLevel:
		return otellog.SeverityDebug
	case level == zapcore.InfoLevel:
		return otellog.SeverityInfo
	case level == zapcore.WarnLevel:
		return otellog.SeverityWarn
	case level >= zapcore.ErrorLevel && level < zapcore.DPanicLevel:
		return otellog.SeverityError
	default:
		return otellog.SeverityFatal
	}
}

// TextFromZap returns a canonical severity text.
func TextFromZap(level zapcore.Level) string {
	return strings.ToUpper(level.String())
}

// FromSlog converts a slog level to OpenTelemetry severity.
func FromSlog(level slog.Level) otellog.Severity {
	switch {
	case level <= slog.LevelDebug:
		return otellog.SeverityDebug
	case level < slog.LevelWarn:
		return otellog.SeverityInfo
	case level < slog.LevelError:
		return otellog.SeverityWarn
	default:
		return otellog.SeverityError
	}
}

// TextFromSlog returns a canonical severity text.
func TextFromSlog(level slog.Level) string {
	return strings.ToUpper(level.String())
}
