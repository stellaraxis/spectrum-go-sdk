package config

import "testing"

func TestNormalizeDevelopmentDefaults(t *testing.T) {
	cfg, err := Config{
		ServiceName: "user-service",
		Environment: "dev",
	}.Normalize()
	if err != nil {
		t.Fatalf("normalize config: %v", err)
	}

	if cfg.Output != OutputConsole {
		t.Fatalf("expected output %q, got %q", OutputConsole, cfg.Output)
	}
	if cfg.Format != FormatConsole {
		t.Fatalf("expected format %q, got %q", FormatConsole, cfg.Format)
	}
}

func TestApplyEnv(t *testing.T) {
	t.Setenv(stellarAppName, "billing-service")
	t.Setenv(stellarEnv, "prod")
	t.Setenv(spectrumEndpoint, "localhost:4317")

	cfg := Default()
	if err := cfg.ApplyEnv(); err != nil {
		t.Fatalf("apply env: %v", err)
	}

	if cfg.ServiceName != "billing-service" {
		t.Fatalf("unexpected service name %q", cfg.ServiceName)
	}
	if cfg.Environment != "prod" {
		t.Fatalf("unexpected environment %q", cfg.Environment)
	}
	if cfg.Endpoint != "localhost:4317" {
		t.Fatalf("unexpected endpoint %q", cfg.Endpoint)
	}
}

func TestNormalizeUsesStellarMetadataDefaults(t *testing.T) {
	t.Setenv(stellarAppName, "order-service")
	t.Setenv(stellarAppNamespace, "stellar.trade")
	t.Setenv(stellarAppVersion, "2.3.1")
	t.Setenv(stellarAppInstanceID, "order-service-0")
	t.Setenv(stellarEnv, "prod")
	t.Setenv(stellarCluster, "cluster-hz-prod-01")
	t.Setenv(stellarRegion, "cn-east-1")
	t.Setenv(stellarZone, "cn-east-1a")

	cfg, err := Config{}.Normalize()
	if err != nil {
		t.Fatalf("normalize config: %v", err)
	}

	if cfg.ServiceName != "order-service" {
		t.Fatalf("unexpected service name %q", cfg.ServiceName)
	}
	if cfg.ServiceNamespace != "stellar.trade" {
		t.Fatalf("unexpected service namespace %q", cfg.ServiceNamespace)
	}
	if cfg.ServiceVersion != "2.3.1" {
		t.Fatalf("unexpected service version %q", cfg.ServiceVersion)
	}
	if cfg.ServiceInstanceID != "order-service-0" {
		t.Fatalf("unexpected service instance id %q", cfg.ServiceInstanceID)
	}
	if cfg.Cluster != "cluster-hz-prod-01" {
		t.Fatalf("unexpected cluster %q", cfg.Cluster)
	}
}

func TestSpectrumEnvOverridesStellarEnv(t *testing.T) {
	t.Setenv(stellarAppName, "order-service")
	t.Setenv(spectrumServiceName, "spectrum-order-service")

	cfg := Default()
	if err := cfg.ApplyEnv(); err != nil {
		t.Fatalf("apply env: %v", err)
	}

	if cfg.ServiceName != "spectrum-order-service" {
		t.Fatalf("unexpected service name %q", cfg.ServiceName)
	}
}
