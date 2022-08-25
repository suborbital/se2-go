package compute

import "net/url"

// Config combines the common configuration options for the three
// Suborbital Compute APIs (Administrative, Builder, and Execution)
type Config struct {
	executionURL *url.URL
	adminURL     *url.URL
	builderURL   *url.URL
}

// DefaultConfig takes the given host and creates a Compute config with the K8S default ports.
// Everything except the scheme and hostname are considered. You need to provide your builder
// host domain.
func DefaultConfig(builderHost string) (*Config, error) {
	execUrl, err := url.Parse("http://scc-atmo-service.suborbital.svc.cluster.local:80")
	if err != nil {
		return nil, err
	}

	adminUrl, err := url.Parse("http://scc-controlplane-service.suborbital.svc.cluster.local:8081")
	if err != nil {
		return nil, err
	}

	builderUrl, err := url.Parse(builderHost)
	if err != nil {
		return nil, err
	}
	builderUrl.Scheme = "https"

	conf := &Config{
		executionURL: execUrl,
		adminURL:     adminUrl,
		builderURL:   builderUrl,
	}

	return conf, nil
}

// LocalConfig generates a Configuration for Compute running in docker-compose
func LocalConfig() *Config {
	conf := &Config{
		executionURL: &url.URL{
			Scheme: "http",
			Host:   "local.suborbital.network:8080",
		},
		adminURL: &url.URL{
			Scheme: "http",
			Host:   "local.suborbital.network:8081",
		},
		builderURL: &url.URL{
			Scheme: "http",
			Host:   "local.suborbital.network:8082",
		},
	}

	return conf
}

// Custom Configuration
func CustomConfig(execHost string, adminHost string, builderHost string) (*Config, error) {
	execUrl, err := url.Parse(execHost)
	if err != nil {
		return nil, err
	}

	adminUrl, err := url.Parse(adminHost)
	if err != nil {
		return nil, err
	}

	builderUrl, err := url.Parse(builderHost)
	if err != nil {
		return nil, err
	}

	conf := &Config{
		executionURL: execUrl,
		adminURL:     adminUrl,
		builderURL:   builderUrl,
	}

	return conf, nil
}
