package se2

import "net/url"

// Config combines the common configuration options for the three
// Suborbital Extension Engine APIs (Administrative, Builder, and Execution)
type Config struct {
	ExecutionURL *url.URL
	AdminURL     *url.URL
	BuilderURL   *url.URL
}

// DefaultConfig takes the given host and creates a SE2 config with the K8S default ports.
// Everything except the scheme and hostname are considered. You need to provide your builder
// host domain.
func DefaultConfig(builderHost string) (*Config, error) {
	execUrl, err := url.Parse("http://e2core-service.suborbital.svc.cluster.local:80")
	if err != nil {
		return nil, err
	}

	adminUrl, err := url.Parse("http://se2-controlplane-service.suborbital.svc.cluster.local:8081")
	if err != nil {
		return nil, err
	}

	builderUrl, err := url.Parse(builderHost)
	if err != nil {
		return nil, err
	}
	builderUrl.Scheme = "https"

	conf := &Config{
		ExecutionURL: execUrl,
		AdminURL:     adminUrl,
		BuilderURL:   builderUrl,
	}

	return conf, nil
}

// LocalConfig generates a Configuration for SE2 running in docker-compose
func LocalConfig() *Config {
	conf := &Config{
		ExecutionURL: &url.URL{
			Scheme: "http",
			Host:   "local.suborbital.network:8080",
		},
		AdminURL: &url.URL{
			Scheme: "http",
			Host:   "local.suborbital.network:8081",
		},
		BuilderURL: &url.URL{
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
		ExecutionURL: execUrl,
		AdminURL:     adminUrl,
		BuilderURL:   builderUrl,
	}

	return conf, nil
}
