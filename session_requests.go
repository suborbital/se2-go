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

type CreateSessionRequest struct {
	Plugin    string `json:"fn"`
	Namespace string `json:"namespace"`
}

type CreateSessionResponse struct {
	Token string `json:"token"`
}

// CreateSession will create a session for a given tenant, namespace, and plugin to be used in the builder.
func (c *Client2) CreateSession(ctx context.Context, tenantName, namespace, plugin string) (CreateSessionResponse, error) {
	// Check arguments.
	if tenantName == "" {
		return CreateSessionResponse{}, errors.New("tenant name cannot be blank")
	}

	if namespace == "" {
		return CreateSessionResponse{}, errors.New("namespace cannot be blank")
	}

	if plugin == "" {
		return CreateSessionResponse{}, errors.New("plugin cannot be blank")
	}

	// Build a body, Dr. Frankenstein!
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(CreateSessionRequest{
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
