package se2

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	modeUnset ServerMode = iota
	ModeStaging
	ModeProduction
	hostProduction     string = "https://api.suborbital.network"
	hostStaging        string = "https://stg.api.suborbital.network"
	hostExecProduction string = "https://edge.suborbital.network"
	hostExecStaging    string = "https://stg.edge.suborbital.network"
	minAccessKeyLength        = 60
)

var (
	ErrNoAccessKey = errors.New("No access key provided, or it's likely malformed.")
	ErrUnknownMode = errors.New("Unknown client mode set. Use one of the ModeStaging or ModeProduction constants.")
)

type ServerMode int

type accessKey struct {
	Key    int    `json:"key"`
	Secret string `json:"secret"`
}

type Client struct {
	httpClient *http.Client
	host       string
	execHost   string
	token      string
}

type ClientOption func(*Client)

// NewClient returns a configured instance of a configured client for SE2. Required parameters are the mode to specify
// whether it's the production or the staging environment, and an access key you can grab from the SE2 admin area for
// an environment.
//
// By default, the underlying http client has a 60-second timeout. Otherwise, you can use the
// WithHttpClient(*http.Client) function to use your own configured version for it.
func NewClient(mode ServerMode, token string, options ...ClientOption) (*Client, error) {
	// Create zero value client with default http client.
	nc := Client{
		httpClient: defaultHttpClient(),
	}

	// Set hosts based on mode, if somehow it's an unknown mode, return error.
	switch mode {
	case ModeStaging:
		nc.host = hostStaging
		nc.execHost = hostExecStaging
	case ModeProduction:
		nc.host = hostProduction
		nc.execHost = hostExecProduction
	default:
		return nil, ErrUnknownMode
	}

	// If access key is too short after modifiers, return error.
	if len(token) < minAccessKeyLength {
		return nil, ErrNoAccessKey
	}

	// If access key is not a base64 encoded string, return error.
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, ErrNoAccessKey
	}

	// If access key is not a valid JSON base64 encoded, return error.
	var akUnmarshaled accessKey
	err = json.Unmarshal(decoded, &akUnmarshaled)
	if err != nil {
		return nil, ErrNoAccessKey
	}

	// Save the good token to the client.
	nc.token = token

	// Apply all the modifiers.
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

// WithHttpClient allows you to configure the http.Client used in the SE2 client.
func WithHttpClient(client *http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = client
	}
}

// do is the meat of the client, every other exported method sets up the request and the context.
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.token)

	return c.httpClient.Do(req)
}
