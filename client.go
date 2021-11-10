package compute

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	config *Config
}

func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, errors.New("failed to NewClient: config cannot be nil")
	}

	client := &Client{
		config: config,
	}

	return client, nil
}

func NewLocalClient() (*Client, error) {
	return NewClient(LocalConfig())
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
		return res, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return res, errors.Errorf("API returned non-successful status: %s", res.Status)
	}

	return res, nil
}
