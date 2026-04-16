package config

import "time"

const (
	DefaultProtocol = "grpc"
	DefaultFormat   = "json"
	DefaultOutput   = "otlp"

	OutputOTLP    = "otlp"
	OutputStdout  = "stdout"
	OutputStderr  = "stderr"
	OutputConsole = "console"

	FormatJSON    = "json"
	FormatConsole = "console"
)

// Config defines the runtime configuration used by the SDK.
type Config struct {
	ServiceName        string
	ServiceNamespace   string
	ServiceVersion     string
	ServiceInstanceID  string
	Environment        string
	Cluster            string
	Region             string
	Zone               string
	IDC                string
	HostName           string
	HostIP             string
	NodeName           string
	K8sNamespace       string
	PodName            string
	PodIP              string
	ContainerName      string
	Endpoint           string
	Insecure           bool
	Protocol           string
	Format             string
	Output             string
	Level              string
	Development        bool
	EnableCaller       bool
	EnableStacktrace   bool
	BatchTimeout       time.Duration
	ExportTimeout      time.Duration
	MaxBatchSize       int
	MaxQueueSize       int
	Headers            map[string]string
	ResourceAttributes map[string]string
}

// Default returns a baseline configuration with sensible defaults.
func Default() Config {
	return Config{
		Protocol:         DefaultProtocol,
		Format:           DefaultFormat,
		Output:           DefaultOutput,
		Level:            "info",
		Insecure:         true,
		BatchTimeout:     5 * time.Second,
		ExportTimeout:    3 * time.Second,
		MaxBatchSize:     512,
		MaxQueueSize:     2048,
		EnableCaller:     true,
		EnableStacktrace: true,
	}
}

// LoadFromEnv returns a config populated from environment variables.
func LoadFromEnv() (Config, error) {
	cfg := Default()
	if err := cfg.ApplyEnv(); err != nil {
		return Config{}, err
	}
	return cfg.Normalize()
}

// IsDevelopment reports whether the config should run in development mode.
func (c Config) IsDevelopment() bool {
	if c.Development {
		return true
	}
	switch c.Environment {
	case "", "dev", "local", "development":
		return true
	default:
		return false
	}
}

// Normalize applies derived defaults and validates the configuration.
func (c Config) Normalize() (Config, error) {
	cfg := c
	cfg.applyGlobalMetadataDefaults()

	if cfg.Protocol == "" {
		cfg.Protocol = DefaultProtocol
	}
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.BatchTimeout <= 0 {
		cfg.BatchTimeout = 5 * time.Second
	}
	if cfg.ExportTimeout <= 0 {
		cfg.ExportTimeout = 3 * time.Second
	}
	if cfg.MaxBatchSize <= 0 {
		cfg.MaxBatchSize = 512
	}
	if cfg.MaxQueueSize <= 0 {
		cfg.MaxQueueSize = 2048
	}
	if cfg.Output == "" {
		if cfg.IsDevelopment() {
			cfg.Output = OutputConsole
		} else {
			cfg.Output = OutputOTLP
		}
	}
	if cfg.Format == "" {
		if cfg.Output == OutputConsole {
			cfg.Format = FormatConsole
		} else {
			cfg.Format = FormatJSON
		}
	}
	if cfg.IsDevelopment() && cfg.Output == OutputOTLP && cfg.Endpoint == "" {
		cfg.Endpoint = "localhost:4317"
	}
	if !cfg.IsDevelopment() && cfg.Endpoint == "" {
		cfg.Endpoint = "localhost:4317"
	}
	if cfg.Headers == nil {
		cfg.Headers = map[string]string{}
	}
	if cfg.ResourceAttributes == nil {
		cfg.ResourceAttributes = map[string]string{}
	}

	return cfg, cfg.Validate()
}
