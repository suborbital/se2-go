package se2

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Client is used for interacting with the Suborbital Compute API
type Client struct {
	config   *Config
	envToken string
}

// NewClient creates a Client with a Config
func NewClient(config *Config, envToken string) (*Client, error) {
	if config == nil {
		return nil, errors.New("failed to NewClient: config cannot be nil")
	}

	client := &Client{
		config:   config,
		envToken: envToken,
	}

	return client, nil
}

// NewLocalClient quickly sets up a Client with a LocalConfig. Useful for testing.
func NewLocalClient(envToken string) (*Client, error) {
	return NewClient(LocalConfig(), envToken)
}

func (c *Client) adminRequestBuilder(method string, endpoint string, body io.Reader) (*http.Request, error) {
	url := *c.config.adminURL
	url.Path = endpoint

	return http.NewRequest(method, url.String(), body)
}

func (c *Client) execRequestBuilder(method string, endpoint string, body io.Reader) (*http.Request, error) {
	url := *c.config.executionURL
	url.Path = endpoint

	return http.NewRequest(method, url.String(), body)
}

func (c *Client) builderRequestBuilder(method string, endpoint string, body io.Reader) (*http.Request, error) {
	url := *c.config.builderURL
	url.Path = endpoint

	return http.NewRequest(method, url.String(), body)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		res.Body.Close()
		return nil, errors.Errorf("API returned non-successful status: %s", res.Status)
	}

	return res, nil
}
