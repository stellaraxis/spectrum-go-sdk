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

	stellspecServiceName      = "STELLSPEC_SERVICE_NAME"
	stellspecServiceNamespace = "STELLSPEC_SERVICE_NAMESPACE"
	stellspecServiceVersion   = "STELLSPEC_SERVICE_VERSION"
	stellspecServiceID        = "STELLSPEC_SERVICE_INSTANCE_ID"
	stellspecEnvironment      = "STELLSPEC_ENVIRONMENT"
	stellspecEndpoint         = "STELLSPEC_ENDPOINT"
	stellspecProtocol         = "STELLSPEC_PROTOCOL"
	stellspecOutput           = "STELLSPEC_OUTPUT"
	stellspecFormat           = "STELLSPEC_FORMAT"
	stellspecLevel            = "STELLSPEC_LEVEL"
	stellspecInsecure         = "STELLSPEC_INSECURE"
	stellspecDevelopment      = "STELLSPEC_DEVELOPMENT"
	stellspecEnableCaller     = "STELLSPEC_ENABLE_CALLER"
	stellspecEnableStacktrace = "STELLSPEC_ENABLE_STACKTRACE"
	stellspecBatchTimeout     = "STELLSPEC_BATCH_TIMEOUT"
	stellspecExportTimeout    = "STELLSPEC_EXPORT_TIMEOUT"
	stellspecMaxBatchSize     = "STELLSPEC_MAX_BATCH_SIZE"
	stellspecMaxQueueSize     = "STELLSPEC_MAX_QUEUE_SIZE"
	stellspecFallbackFilePath = "STELLSPEC_FALLBACK_FILE_PATH"
	stellspecRetryEnabled     = "STELLSPEC_RETRY_ENABLED"
	stellspecRetryInitial     = "STELLSPEC_RETRY_INITIAL_INTERVAL"
	stellspecRetryMaxInterval = "STELLSPEC_RETRY_MAX_INTERVAL"
	stellspecRetryMaxElapsed  = "STELLSPEC_RETRY_MAX_ELAPSED_TIME"
)

// ApplyEnv loads global Stellar metadata first, then applies Stellspec-specific overrides.
func (c *Config) ApplyEnv() error {
	if err := c.ApplyStellarEnv(); err != nil {
		return err
	}
	return c.ApplyStellspecEnv()
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

// ApplyStellspecEnv loads Stellspec-specific overrides.
func (c *Config) ApplyStellspecEnv() error {
	setString(&c.ServiceName, stellspecServiceName)
	setString(&c.ServiceNamespace, stellspecServiceNamespace)
	setString(&c.ServiceVersion, stellspecServiceVersion)
	setString(&c.ServiceInstanceID, stellspecServiceID)
	setString(&c.Environment, stellspecEnvironment)
	setString(&c.Endpoint, stellspecEndpoint)
	setString(&c.Protocol, stellspecProtocol)
	setString(&c.Output, stellspecOutput)
	setString(&c.Format, stellspecFormat)
	setString(&c.Level, stellspecLevel)
	setString(&c.FallbackFilePath, stellspecFallbackFilePath)

	if err := setBool(&c.Insecure, stellspecInsecure); err != nil {
		return err
	}
	if err := setBool(&c.Development, stellspecDevelopment); err != nil {
		return err
	}
	if err := setBool(&c.EnableCaller, stellspecEnableCaller); err != nil {
		return err
	}
	if err := setBool(&c.EnableStacktrace, stellspecEnableStacktrace); err != nil {
		return err
	}
	if err := setDuration(&c.BatchTimeout, stellspecBatchTimeout); err != nil {
		return err
	}
	if err := setDuration(&c.ExportTimeout, stellspecExportTimeout); err != nil {
		return err
	}
	if err := setInt(&c.MaxBatchSize, stellspecMaxBatchSize); err != nil {
		return err
	}
	if err := setInt(&c.MaxQueueSize, stellspecMaxQueueSize); err != nil {
		return err
	}
	if err := setRetryBool(&c.Retry.Enabled, stellspecRetryEnabled); err != nil {
		return err
	}
	if err := setRetryDuration(&c.Retry.InitialInterval, stellspecRetryInitial); err != nil {
		return err
	}
	if err := setRetryDuration(&c.Retry.MaxInterval, stellspecRetryMaxInterval); err != nil {
		return err
	}
	if err := setRetryDuration(&c.Retry.MaxElapsedTime, stellspecRetryMaxElapsed); err != nil {
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
