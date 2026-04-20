package config

import (
	"fmt"
	"strings"
)

// Validate checks whether the config can be used safely.
func (c Config) Validate() error {
	if strings.TrimSpace(c.ServiceName) == "" {
		return fmt.Errorf("service name is required")
	}

	switch c.Protocol {
	case DefaultProtocol:
	default:
		return fmt.Errorf("unsupported protocol %q", c.Protocol)
	}

	switch c.Output {
	case OutputOTLP, OutputStdout, OutputStderr, OutputConsole:
	default:
		return fmt.Errorf("unsupported output %q", c.Output)
	}

	switch c.Format {
	case FormatJSON, FormatConsole:
	default:
		return fmt.Errorf("unsupported format %q", c.Format)
	}

	if c.Output == OutputOTLP && strings.TrimSpace(c.Endpoint) == "" {
		return fmt.Errorf("endpoint is required when output is %q", OutputOTLP)
	}
	if c.MaxBatchSize <= 0 {
		return fmt.Errorf("max batch size must be greater than zero")
	}
	if c.MaxQueueSize <= 0 {
		return fmt.Errorf("max queue size must be greater than zero")
	}
	if c.Retry.Enabled == nil {
		return fmt.Errorf("retry enabled must not be nil")
	}
	if c.Retry.InitialInterval == nil {
		return fmt.Errorf("retry initial interval must not be nil")
	}
	if c.Retry.MaxInterval == nil {
		return fmt.Errorf("retry max interval must not be nil")
	}
	if c.Retry.MaxElapsedTime == nil {
		return fmt.Errorf("retry max elapsed time must not be nil")
	}
	if *c.Retry.InitialInterval <= 0 {
		return fmt.Errorf("retry initial interval must be greater than zero")
	}
	if *c.Retry.MaxInterval <= 0 {
		return fmt.Errorf("retry max interval must be greater than zero")
	}
	if *c.Retry.MaxElapsedTime < 0 {
		return fmt.Errorf("retry max elapsed time must not be negative")
	}

	return nil
}
