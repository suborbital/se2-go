package compute

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config *Config, httpClient *http.Client) (*Client, error) {
	if config == nil {
		return nil, errors.New("failed to NewClient: config cannot be nil")
	}

	if httpClient == nil {
		return nil, errors.New("failed to NewClient: httpClient cannot be nil")
	}

	client := &Client{
		config:     config,
		httpClient: httpClient,
	}

	return client, nil
}

func NewLocalClient() (*Client, error) {
	return NewClient(LocalConfig(), http.DefaultClient)
}

func (c *Client) NewRunnable(environment, customerID, namespace, fnName string) *Runnable {
	runnable := &Runnable{
		environment:  environment,
		customerID:   customerID,
		namespace:    namespace,
		functionName: fnName,
	}

	return runnable
}

func (c *Client) BuildWith(runnable *Runnable, fn Function) (string, *http.Response, error) {
	return "", nil, nil
}

func (c *Client) adminRequestBuilder(method string, endpoint string, body io.Reader) (*http.Request, error) {
	url := *c.config.adminURL
	url.Path = endpoint

	return http.NewRequest(method, url.String(), body)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)

	if err != nil {
		return res, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return res, errors.Errorf("API returned non-successful status: %s", res.Status)
	}

	return res, nil
}
