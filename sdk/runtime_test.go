package sdk

import (
	"context"
	"io"
	"testing"

	"github.com/stellaraxis/spectrum-go-sdk/config"
	"go.opentelemetry.io/otel/attribute"
)

func TestNewRuntime(t *testing.T) {
	runtime, err := New(context.Background(), config.Config{
		ServiceName: "user-service",
		Environment: "dev",
	}, WithWriters(io.Discard, io.Discard))
	if err != nil {
		t.Fatalf("new runtime: %v", err)
	}

	if runtime.LoggerProvider() == nil {
		t.Fatal("logger provider should not be nil")
	}
	if runtime.Resource() == nil {
		t.Fatal("resource should not be nil")
	}
	if err := runtime.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown runtime: %v", err)
	}
}

func TestNewRuntimeBuildsEnterpriseResourceAttributes(t *testing.T) {
	runtime, err := New(context.Background(), config.Config{
		ServiceName:       "user-service",
		Environment:       "prod",
		ServiceVersion:    "1.0.0",
		ServiceInstanceID: "user-service-0",
		Cluster:           "cluster-sh-prod-01",
		Region:            "cn-east-1",
		Zone:              "cn-east-1a",
		IDC:               "sh-a",
		HostName:          "node-01",
		HostIP:            "10.0.0.10",
		NodeName:          "worker-01",
		K8sNamespace:      "trade",
		PodName:           "user-service-0",
		PodIP:             "10.1.0.20",
		ContainerName:     "app",
	}, WithWriters(io.Discard, io.Discard))
	if err != nil {
		t.Fatalf("new runtime: %v", err)
	}
	t.Cleanup(func() {
		_ = runtime.Shutdown(context.Background())
	})

	attrSet := runtime.Resource().Set()
	assertAttr(t, attrSet, "service.instance.id", "user-service-0")
	assertAttr(t, attrSet, "deployment.environment.name", "prod")
	assertAttr(t, attrSet, "stellar.cluster", "cluster-sh-prod-01")
	assertAttr(t, attrSet, "cloud.region", "cn-east-1")
	assertAttr(t, attrSet, "cloud.availability_zone", "cn-east-1a")
	assertAttr(t, attrSet, "stellar.idc", "sh-a")
	assertAttr(t, attrSet, "host.name", "node-01")
	assertAttr(t, attrSet, "host.ip", "10.0.0.10")
	assertAttr(t, attrSet, "k8s.node.name", "worker-01")
	assertAttr(t, attrSet, "k8s.namespace.name", "trade")
	assertAttr(t, attrSet, "k8s.pod.name", "user-service-0")
	assertAttr(t, attrSet, "k8s.pod.ip", "10.1.0.20")
	assertAttr(t, attrSet, "container.name", "app")
}

func assertAttr(t *testing.T, set *attribute.Set, key string, expected string) {
	t.Helper()

	value, ok := set.Value(attribute.Key(key))
	if !ok {
		t.Fatalf("attribute %q not found", key)
	}
	if value.AsString() != expected {
		t.Fatalf("unexpected value for %q: got %q want %q", key, value.AsString(), expected)
	}
}
