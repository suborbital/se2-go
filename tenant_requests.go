package se2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	pathTenantByName = "/environment/v1/tenant/%s"
)

type GetTenantResponse struct {
	id, environment, name, description string
}

// GetTenantByName returns the tenant
func (c *Client2) GetTenantByName(ctx context.Context, name string) (GetTenantResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.host+fmt.Sprintf(pathTenantByName, name), nil)
	if err != nil {
		return GetTenantResponse{}, errors.Wrap(err, "GetTenantByName: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return GetTenantResponse{}, errors.Wrap(err, "GetTenantByName: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var t GetTenantResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return GetTenantResponse{}, errors.Wrap(err, "GetTenantByName: dec.Decode")
	}

	return t, nil
}

type CreateTenantRequest struct {
	Description string `json:"description"`
}

type CreateTenantResponse struct {
	AuthorizedParty string `json:"authorized_party"`
	Environment     string `json:"environment"`
	ID              string `json:"id"`
}

func (c *Client2) CreateTenantByName(ctx context.Context, name, description string) (CreateTenantResponse, error) {
	if name == "" {
		return CreateTenantResponse{}, errors.New("CreateTenantByName: tenant name cannot be empty")
	}

	var requestBody io.Reader
	if description != "" {
		m, err := json.Marshal(CreateTenantRequest{Description: description})
		if err != nil {
			return CreateTenantResponse{}, errors.Wrap(err, "CreateTenantByName: json marshal create tenant request with description")
		}

		requestBody = bytes.NewReader(m)
	}

	req, err := http.NewRequest(http.MethodPost, c.host+fmt.Sprintf(pathTenantByName, name), requestBody)
	if err != nil {
		return CreateTenantResponse{}, errors.Wrap(err, "CreateTenantByName: http.NewRequest for POST create tenant")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return CreateTenantResponse{}, errors.Wrap(err, "CreateTenantByName: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var t CreateTenantResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return CreateTenantResponse{}, errors.Wrap(err, "CreateTenantByName: dec.Decode")
	}

	return t, nil
}
