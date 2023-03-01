package se2

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client is used for interacting with the SE2 API
type Client struct {
	config   *Config
	envToken string
}

const (
	Production         ServerURL = "https://api.suborbital.network"
	Staging            ServerURL = "https://stg.api.suborbital.network"
	minAccessKeyLength           = 60
)

var (
	ErrNoAccessKey = errors.New("No access key provided, or it's likely malformed")
)

type ServerURL string
type accessKey struct {
	Key    int    `json:"key"`
	Secret string `json:"secret"`
}

type Client2 struct {
	httpClient *http.Client
	host       ServerURL
}

type ClientOption func(*Client2)

func NewClient2(host ServerURL, ak string, options ...ClientOption) (*Client2, error) {
	fmt.Printf("len ak: %d\n", len(ak))
	if len(ak) < minAccessKeyLength {
		fmt.Printf("returning because it's too bla")
		return nil, ErrNoAccessKey
	}

	decoded, err := base64.StdEncoding.DecodeString(ak)
	if err != nil {
		return nil, ErrNoAccessKey
	}

	var akUnmarshaled accessKey
	err = json.Unmarshal(decoded, &akUnmarshaled)
	if err != nil {
		return nil, errors.Wrap(err, ErrNoAccessKey.Error())
	}

	nc := Client2{
		httpClient: defaultHttpClient(),
		host:       host,
	}

	for _, o := range options {
		o(&nc)
	}

	return &nc, nil
}

func defaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
	}
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
