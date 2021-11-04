package compute

import (
	admin "github.com/suborbital/compute-go/compute/administrative"
	builder "github.com/suborbital/compute-go/compute/builder"
	execution "github.com/suborbital/compute-go/compute/execution"
)

// Config combines the common configuration options for the three
// Suborbital Compute APIs (Administrative, Builder, and Execution)
type Config struct {
	adminConfig     *admin.Configuration
	builderConfig   *builder.Configuration
	executionConfig *execution.Configuration

	// TODO: pass variables into each config (host, port, etc.)
}

func DefaultConfig() *Config {
	conf := &Config{
		adminConfig:     admin.NewConfiguration(),
		builderConfig:   builder.NewConfiguration(),
		executionConfig: execution.NewConfiguration(),
	}

	return conf
}
