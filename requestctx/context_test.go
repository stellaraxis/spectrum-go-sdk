package requestctx

import (
	"context"
	"testing"
)

func TestWithValuesAndFromContext(t *testing.T) {
	ctx := WithValues(context.Background(), Values{
		RequestID:   "req-1001",
		TenantID:    "tenant-a",
		TraceParent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	})

	values := FromContext(ctx)
	if values.RequestID != "req-1001" {
		t.Fatalf("unexpected request id %q", values.RequestID)
	}
	if values.TenantID != "tenant-a" {
		t.Fatalf("unexpected tenant id %q", values.TenantID)
	}
	if values.TraceParent == "" {
		t.Fatal("traceparent should not be empty")
	}
}

func TestFields(t *testing.T) {
	ctx := WithValues(context.Background(), Values{
		RequestID:    "req-2001",
		TenantID:     "tenant-b",
		SourceApp:    "event-horizon-external",
		TraceParent:  "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		TraceState:   "vendor=value",
		Environment:  "prod",
		ForwardedFor: "10.0.0.1",
	})

	fields := Fields(ctx)
	if fields["request_id"] != "req-2001" {
		t.Fatalf("unexpected request_id %q", fields["request_id"])
	}
	if fields["tenant_id"] != "tenant-b" {
		t.Fatalf("unexpected tenant_id %q", fields["tenant_id"])
	}
	if fields["source_app"] != "event-horizon-external" {
		t.Fatalf("unexpected source_app %q", fields["source_app"])
	}
	if fields["traceparent"] == "" {
		t.Fatal("traceparent should not be empty")
	}
}
