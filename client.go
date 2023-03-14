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
	defaultTimeout            = 60 * time.Second
	emptyString        string = ""
)

var (
	ErrNoAccessKey = errors.New("No access key provided, or it's likely malformed.")
	ErrUnknownMode = errors.New("Unknown client mode set. Use one of the ModeStaging or ModeProduction constants.")
)

// ServerMode is an alias type to help ensure that only the options we declared here can be used.
type ServerMode int

// accessKey is a transitive type that helps the NewClient constructor make sure that the passed in token is of the
// correct form and structure.
type accessKey struct {
	Key    int    `json:"key"`
	Secret string `json:"secret"`
}

// Client holds our configured http client, its methods, and all the functionality necessary so users can interact with
// our API without having to call actual http requests.
type Client struct {
	httpClient *http.Client
	host       string
	execHost   string
	token      string
}

// ClientOption is a function signature users can use to configure different parts of the client. They are run at the
// very end of initialization, once we know that the mode and the access key are correct.
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

// defaultHttpClient returns an http.Client with a 60 second timeout that's used until the users decide to change it by
// use the WithHttpClient function.
func defaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: defaultTimeout,
	}
}

// WithHttpClient allows you to configure the http.Client used in the SE2 client.
func WithHttpClient(client *http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = client
	}
}

// do is the meat of the client, every other admin level exported method uses this. Its main job is to attach the
// context and the access key to outgoing requests.
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.token)

	return c.httpClient.Do(req)
}

// sessionDo is a common method to work with requests against the builder where a session token is needed instead of the
// environment token that the do method uses.
func (c *Client) sessionDo(ctx context.Context, req *http.Request, token CreateSessionResponse) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Add("Authorization", "Bearer "+token.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "c.sessionDo: c.httpClient.Do")
	}

	return res, nil
}
