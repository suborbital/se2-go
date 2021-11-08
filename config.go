package compute

import "net/url"

// Config combines the common configuration options for the three
// Suborbital Compute APIs (Administrative, Builder, and Execution)
type Config struct {
	executionURL *url.URL
	adminURL     *url.URL
	builderURL   *url.URL
}

// DefaultConfig takes the given host and creates a Compute config with the default ports.
// Everything except the scheme and hostname are considered.
func DefaultConfig(host string) (*Config, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	hostname := u.Hostname()
	scheme := u.Scheme

	conf := &Config{
		executionURL: &url.URL{
			Scheme: scheme,
			Host:   hostname + ":8080",
		},
		adminURL: &url.URL{
			Scheme: scheme,
			Host:   hostname + ":8081",
		},
		builderURL: &url.URL{
			Scheme: scheme,
			Host:   hostname + ":8082",
		},
	}

	return conf, nil
}

// LocalConfig generates a DefaultConfig() for localhost
func LocalConfig() *Config {
	conf, _ := DefaultConfig("http://local.suborbital.network")
	return conf
}
