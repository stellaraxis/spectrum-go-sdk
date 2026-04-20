package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	stellarAppName       = "STELLAR_APP_NAME"
	stellarAppNamespace  = "STELLAR_APP_NAMESPACE"
	stellarAppVersion    = "STELLAR_APP_VERSION"
	stellarAppInstanceID = "STELLAR_APP_INSTANCE_ID"
	stellarEnv           = "STELLAR_ENV"
	stellarCluster       = "STELLAR_CLUSTER"
	stellarRegion        = "STELLAR_REGION"
	stellarZone          = "STELLAR_ZONE"
	stellarIDC           = "STELLAR_IDC"
	stellarHostName      = "STELLAR_HOST_NAME"
	stellarHostIP        = "STELLAR_HOST_IP"
	stellarNodeName      = "STELLAR_NODE_NAME"
	stellarK8sNamespace  = "STELLAR_K8S_NAMESPACE"
	stellarPodName       = "STELLAR_POD_NAME"
	stellarPodIP         = "STELLAR_POD_IP"
	stellarContainerName = "STELLAR_CONTAINER_NAME"

	spectrumServiceName      = "SPECTRUM_SERVICE_NAME"
	spectrumServiceNamespace = "SPECTRUM_SERVICE_NAMESPACE"
	spectrumServiceVersion   = "SPECTRUM_SERVICE_VERSION"
	spectrumServiceID        = "SPECTRUM_SERVICE_INSTANCE_ID"
	spectrumEnvironment      = "SPECTRUM_ENVIRONMENT"
	spectrumEndpoint         = "SPECTRUM_ENDPOINT"
	spectrumProtocol         = "SPECTRUM_PROTOCOL"
	spectrumOutput           = "SPECTRUM_OUTPUT"
	spectrumFormat           = "SPECTRUM_FORMAT"
	spectrumLevel            = "SPECTRUM_LEVEL"
	spectrumInsecure         = "SPECTRUM_INSECURE"
	spectrumDevelopment      = "SPECTRUM_DEVELOPMENT"
	spectrumEnableCaller     = "SPECTRUM_ENABLE_CALLER"
	spectrumEnableStacktrace = "SPECTRUM_ENABLE_STACKTRACE"
	spectrumBatchTimeout     = "SPECTRUM_BATCH_TIMEOUT"
	spectrumExportTimeout    = "SPECTRUM_EXPORT_TIMEOUT"
	spectrumMaxBatchSize     = "SPECTRUM_MAX_BATCH_SIZE"
	spectrumMaxQueueSize     = "SPECTRUM_MAX_QUEUE_SIZE"
	spectrumFallbackFilePath = "SPECTRUM_FALLBACK_FILE_PATH"
	spectrumRetryEnabled     = "SPECTRUM_RETRY_ENABLED"
	spectrumRetryInitial     = "SPECTRUM_RETRY_INITIAL_INTERVAL"
	spectrumRetryMaxInterval = "SPECTRUM_RETRY_MAX_INTERVAL"
	spectrumRetryMaxElapsed  = "SPECTRUM_RETRY_MAX_ELAPSED_TIME"
)

// ApplyEnv loads global Stellar metadata first, then applies Spectrum-specific overrides.
func (c *Config) ApplyEnv() error {
	if err := c.ApplyStellarEnv(); err != nil {
		return err
	}
	return c.ApplySpectrumEnv()
}

// ApplyStellarEnv loads global enterprise application metadata.
func (c *Config) ApplyStellarEnv() error {
	setString(&c.ServiceName, stellarAppName)
	setString(&c.ServiceNamespace, stellarAppNamespace)
	setString(&c.ServiceVersion, stellarAppVersion)
	setString(&c.ServiceInstanceID, stellarAppInstanceID)
	setString(&c.Environment, stellarEnv)
	setString(&c.Cluster, stellarCluster)
	setString(&c.Region, stellarRegion)
	setString(&c.Zone, stellarZone)
	setString(&c.IDC, stellarIDC)
	setString(&c.HostName, stellarHostName)
	setString(&c.HostIP, stellarHostIP)
	setString(&c.NodeName, stellarNodeName)
	setString(&c.K8sNamespace, stellarK8sNamespace)
	setString(&c.PodName, stellarPodName)
	setString(&c.PodIP, stellarPodIP)
	setString(&c.ContainerName, stellarContainerName)
	return nil
}

// ApplySpectrumEnv loads Spectrum-specific overrides.
func (c *Config) ApplySpectrumEnv() error {
	setString(&c.ServiceName, spectrumServiceName)
	setString(&c.ServiceNamespace, spectrumServiceNamespace)
	setString(&c.ServiceVersion, spectrumServiceVersion)
	setString(&c.ServiceInstanceID, spectrumServiceID)
	setString(&c.Environment, spectrumEnvironment)
	setString(&c.Endpoint, spectrumEndpoint)
	setString(&c.Protocol, spectrumProtocol)
	setString(&c.Output, spectrumOutput)
	setString(&c.Format, spectrumFormat)
	setString(&c.Level, spectrumLevel)
	setString(&c.FallbackFilePath, spectrumFallbackFilePath)

	if err := setBool(&c.Insecure, spectrumInsecure); err != nil {
		return err
	}
	if err := setBool(&c.Development, spectrumDevelopment); err != nil {
		return err
	}
	if err := setBool(&c.EnableCaller, spectrumEnableCaller); err != nil {
		return err
	}
	if err := setBool(&c.EnableStacktrace, spectrumEnableStacktrace); err != nil {
		return err
	}
	if err := setDuration(&c.BatchTimeout, spectrumBatchTimeout); err != nil {
		return err
	}
	if err := setDuration(&c.ExportTimeout, spectrumExportTimeout); err != nil {
		return err
	}
	if err := setInt(&c.MaxBatchSize, spectrumMaxBatchSize); err != nil {
		return err
	}
	if err := setInt(&c.MaxQueueSize, spectrumMaxQueueSize); err != nil {
		return err
	}
	if err := setRetryBool(&c.Retry.Enabled, spectrumRetryEnabled); err != nil {
		return err
	}
	if err := setRetryDuration(&c.Retry.InitialInterval, spectrumRetryInitial); err != nil {
		return err
	}
	if err := setRetryDuration(&c.Retry.MaxInterval, spectrumRetryMaxInterval); err != nil {
		return err
	}
	if err := setRetryDuration(&c.Retry.MaxElapsedTime, spectrumRetryMaxElapsed); err != nil {
		return err
	}

	return nil
}

func (c *Config) applyGlobalMetadataDefaults() {
	fillStringIfEmpty(&c.ServiceName, stellarAppName)
	fillStringIfEmpty(&c.ServiceNamespace, stellarAppNamespace)
	fillStringIfEmpty(&c.ServiceVersion, stellarAppVersion)
	fillStringIfEmpty(&c.ServiceInstanceID, stellarAppInstanceID)
	fillStringIfEmpty(&c.Environment, stellarEnv)
	fillStringIfEmpty(&c.Cluster, stellarCluster)
	fillStringIfEmpty(&c.Region, stellarRegion)
	fillStringIfEmpty(&c.Zone, stellarZone)
	fillStringIfEmpty(&c.IDC, stellarIDC)
	fillStringIfEmpty(&c.HostName, stellarHostName)
	fillStringIfEmpty(&c.HostIP, stellarHostIP)
	fillStringIfEmpty(&c.NodeName, stellarNodeName)
	fillStringIfEmpty(&c.K8sNamespace, stellarK8sNamespace)
	fillStringIfEmpty(&c.PodName, stellarPodName)
	fillStringIfEmpty(&c.PodIP, stellarPodIP)
	fillStringIfEmpty(&c.ContainerName, stellarContainerName)
}

func setString(target *string, key string) {
	value, ok := os.LookupEnv(key)
	if ok && value != "" {
		*target = value
	}
}

func fillStringIfEmpty(target *string, key string) {
	if *target != "" {
		return
	}
	setString(target, key)
}

func setBool(target *bool, key string) error {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("parse %s: %w", key, err)
	}
	*target = parsed
	return nil
}

func setInt(target *int, key string) error {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("parse %s: %w", key, err)
	}
	*target = parsed
	return nil
}

func setDuration(target *time.Duration, key string) error {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("parse %s: %w", key, err)
	}
	*target = parsed
	return nil
}

func setRetryBool(target **bool, key string) error {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("parse %s: %w", key, err)
	}
	*target = &parsed
	return nil
}

func setRetryDuration(target **time.Duration, key string) error {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("parse %s: %w", key, err)
	}
	*target = &parsed
	return nil
}
