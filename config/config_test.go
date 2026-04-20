package config

import (
	"testing"
	"time"
)

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
	t.Setenv(spectrumFallbackFilePath, "logs/custom-fallback.log")
	t.Setenv(spectrumRetryEnabled, "false")
	t.Setenv(spectrumRetryInitial, "2s")
	t.Setenv(spectrumRetryMaxInterval, "10s")
	t.Setenv(spectrumRetryMaxElapsed, "20s")

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
	if cfg.FallbackFilePath != "logs/custom-fallback.log" {
		t.Fatalf("unexpected fallback file path %q", cfg.FallbackFilePath)
	}
	if cfg.Retry.Enabled == nil || *cfg.Retry.Enabled {
		t.Fatalf("unexpected retry enabled: %v", cfg.Retry.Enabled)
	}
	if cfg.Retry.InitialInterval == nil || *cfg.Retry.InitialInterval != 2*time.Second {
		t.Fatalf("unexpected retry initial interval: %v", cfg.Retry.InitialInterval)
	}
	if cfg.Retry.MaxInterval == nil || *cfg.Retry.MaxInterval != 10*time.Second {
		t.Fatalf("unexpected retry max interval: %v", cfg.Retry.MaxInterval)
	}
	if cfg.Retry.MaxElapsedTime == nil || *cfg.Retry.MaxElapsedTime != 20*time.Second {
		t.Fatalf("unexpected retry max elapsed time: %v", cfg.Retry.MaxElapsedTime)
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

func TestNormalizeUsesDefaultFallbackFilePath(t *testing.T) {
	cfg, err := Config{
		ServiceName:      "user-service",
		Environment:      "prod",
		Output:           OutputOTLP,
		Endpoint:         "localhost:4317",
		FallbackFilePath: "",
	}.Normalize()
	if err != nil {
		t.Fatalf("normalize config: %v", err)
	}

	if cfg.FallbackFilePath != DefaultFallbackFilePath {
		t.Fatalf("unexpected fallback file path %q", cfg.FallbackFilePath)
	}
}

func TestNormalizeUsesDefaultRetryConfig(t *testing.T) {
	cfg, err := Config{
		ServiceName: "user-service",
		Environment: "prod",
		Output:      OutputOTLP,
		Endpoint:    "localhost:4317",
	}.Normalize()
	if err != nil {
		t.Fatalf("normalize config: %v", err)
	}

	if cfg.Retry.Enabled == nil || !*cfg.Retry.Enabled {
		t.Fatalf("unexpected retry enabled: %v", cfg.Retry.Enabled)
	}
	if cfg.Retry.InitialInterval == nil || *cfg.Retry.InitialInterval != defaultRetryInitial {
		t.Fatalf("unexpected retry initial interval: %v", cfg.Retry.InitialInterval)
	}
	if cfg.Retry.MaxInterval == nil || *cfg.Retry.MaxInterval != defaultRetryMaxInterval {
		t.Fatalf("unexpected retry max interval: %v", cfg.Retry.MaxInterval)
	}
	if cfg.Retry.MaxElapsedTime == nil || *cfg.Retry.MaxElapsedTime != defaultRetryMaxElapsed {
		t.Fatalf("unexpected retry max elapsed time: %v", cfg.Retry.MaxElapsedTime)
	}
}
