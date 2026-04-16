package requestctx

import (
	"context"
	"strings"
)

type contextKey struct{}

// Values stores normalized request metadata produced by the upper-level framework.
type Values struct {
	RequestID      string
	SessionID      string
	UserID         string
	TenantID       string
	DeviceID       string
	ClientIP       string
	SourceApp      string
	SourceService  string
	SourceRegion   string
	Environment    string
	GrayTag        string
	CanaryTag      string
	TraceParent    string
	TraceState     string
	Baggage        string
	ForwardedFor   string
	ForwardedHost  string
	ForwardedProto string
	RealIP         string
}

// WithValues attaches request values to the context.
func WithValues(ctx context.Context, values Values) context.Context {
	return context.WithValue(ctx, contextKey{}, values)
}

// FromContext loads request values from the context.
func FromContext(ctx context.Context) Values {
	if ctx == nil {
		return Values{}
	}

	values, ok := ctx.Value(contextKey{}).(Values)
	if !ok {
		return Values{}
	}
	return values
}

// Fields returns normalized log fields for the current request context.
func Fields(ctx context.Context) map[string]string {
	values := FromContext(ctx)
	fields := map[string]string{}

	putString(fields, "request_id", values.RequestID)
	putString(fields, "session_id", values.SessionID)
	putString(fields, "user_id", values.UserID)
	putString(fields, "tenant_id", values.TenantID)
	putString(fields, "device_id", values.DeviceID)
	putString(fields, "client_ip", values.ClientIP)
	putString(fields, "source_app", values.SourceApp)
	putString(fields, "source_service", values.SourceService)
	putString(fields, "source_region", values.SourceRegion)
	putString(fields, "env", values.Environment)
	putString(fields, "gray_tag", values.GrayTag)
	putString(fields, "canary_tag", values.CanaryTag)
	putString(fields, "traceparent", values.TraceParent)
	putString(fields, "tracestate", values.TraceState)

	return fields
}

func putString(target map[string]string, key string, value string) {
	if strings.TrimSpace(value) != "" {
		target[key] = value
	}
}
