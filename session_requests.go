package se2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const pathCreateTenantSession = pathTenantByName + "/session"

// createSessionRequest is an internal struct to help with converting data into a json payload we can send against the
// API.
type createSessionRequest struct {
	Plugin    string `json:"fn"`
	Namespace string `json:"namespace"`
}

// CreateSessionResponse has a token inside of it. This token is used in queries against the builder service. Those
// methods will require one of their parameters to be of this type.
type CreateSessionResponse struct {
	Token string `json:"token"`
}

// CreateSession will create a session for a given tenant, namespace, and plugin to be used in the builder. You should
// keep track of the return argument and reuse it in later requests.
func (c *Client) CreateSession(ctx context.Context, tenantName, namespace, plugin string) (CreateSessionResponse, error) {
	// Check arguments.
	if tenantName == emptyString {
		return CreateSessionResponse{}, errors.New("tenant name cannot be blank")
	}

	if namespace == emptyString {
		return CreateSessionResponse{}, errors.New("namespace cannot be blank")
	}

	if plugin == emptyString {
		return CreateSessionResponse{}, errors.New("plugin cannot be blank")
	}

	// Build a body, Dr. Frankenstein!
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(createSessionRequest{
		Plugin:    plugin,
		Namespace: namespace,
	})
	if err != nil {
		return CreateSessionResponse{}, errors.Wrap(err, "CreateSession: json.NewEncoder().Encode")
	}

	// Create the request with the body.
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(c.host+pathCreateTenantSession, tenantName), &body)
	if err != nil {
		return CreateSessionResponse{}, errors.Wrap(err, "CreateSession: http.NewRequest")
	}

	// Do the request.
	res, err := c.do(ctx, req)
	if err != nil {
		return CreateSessionResponse{}, errors.Wrap(err, "CreateSession: c.do")
	}

	// Check response code.
	if res.StatusCode != http.StatusCreated {
		return CreateSessionResponse{}, fmt.Errorf("")
	}

	// Marshal response body into what we need to give back.
	var t CreateSessionResponse
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return CreateSessionResponse{}, errors.Wrap(err, "GetPlugins: dec.Decode")
	}

	return t, nil
}
