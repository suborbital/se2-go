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
	pathTenant       = "/environment/v1/tenant"
	pathTenantByName = pathTenant + "/%s"
)

// TenantResponse captures the JSON data returned from the endpoints.
type TenantResponse struct {
	AuthorizedParty string `json:"authorized_party"`
	Id              string `json:"id"`
	Environment     string `json:"environment"`
	Name            string `json:"name"`
	Description     string `json:"description"`
}

// GetTenantByName returns the tenant by name.
func (c *Client2) GetTenantByName(ctx context.Context, name string) (TenantResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.host+fmt.Sprintf(pathTenantByName, name), nil)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "GetTenantByName: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "GetTenantByName: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var t TenantResponse
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "GetTenantByName: dec.Decode")
	}

	return t, nil
}

// createTenantRequest is used to format the incoming description into a JSON body. Users of this client library should
// not need to use this struct directly.
type createTenantRequest struct {
	Description string `json:"description"`
}

// CreateTenant creates a new tenant with the given name and description. Description is optional, it can be an
// empty string.
func (c *Client2) CreateTenant(ctx context.Context, name, description string) (TenantResponse, error) {
	if name == "" {
		return TenantResponse{}, errors.New("CreateTenant: tenant name cannot be empty")
	}

	var requestBody io.Reader
	if description != "" {
		m, err := json.Marshal(createTenantRequest{Description: description})
		if err != nil {
			return TenantResponse{}, errors.Wrap(err, "CreateTenant: json marshal create tenant request with description")
		}

		requestBody = bytes.NewReader(m)
	}

	req, err := http.NewRequest(http.MethodPost, c.host+fmt.Sprintf(pathTenantByName, name), requestBody)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "CreateTenant: http.NewRequest for POST create tenant")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "CreateTenant: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var t TenantResponse
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "CreateTenant: dec.Decode")
	}

	return t, nil
}

// ListTenantResponse is the unmarshaled response from the endpoint.
type ListTenantResponse struct {
	Tenants []TenantResponse
}

// ListTenants will list the tenants that the configured API key can access.
func (c *Client2) ListTenants(ctx context.Context) (ListTenantResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.host+pathTenant, nil)
	if err != nil {
		return ListTenantResponse{}, errors.Wrap(err, "ListTenants: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return ListTenantResponse{}, errors.Wrap(err, "ListTenants: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var t ListTenantResponse
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return ListTenantResponse{}, errors.Wrap(err, "ListTenants: dec.Decode")
	}

	return t, nil
}

// updateTenantRequest is a struct to help constrain the data and marshal the json based on it that ends up in a request
// body. It's internal only, users of the client do not need to know about the existence of this struct.
type updateTenantRequest struct {
	Description string `json:"description"`
}

// UpdateTenantByName updates the description of the tenant identified by its name. A tenant's name cannot be changed.
func (c *Client2) UpdateTenantByName(ctx context.Context, name, description string) (TenantResponse, error) {
	if name == "" {
		return TenantResponse{}, errors.New("UpdateTenantByName: tenant name cannot be empty")
	}

	m, err := json.Marshal(updateTenantRequest{Description: description})
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "UpdateTenantByName: json marshal update tenant request")
	}

	req, err := http.NewRequest(http.MethodPatch, c.host+fmt.Sprintf(pathTenantByName, name), bytes.NewReader(m))
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "UpdateTenantByName: http.NewRequest for POST create tenant")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "UpdateTenantByName: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var t TenantResponse
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "UpdateTenantByName: dec.Decode")
	}

	return t, nil
}

// DeleteTenantByName deletes the tenant identified by its name.
func (c *Client2) DeleteTenantByName(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("DeleteTenantByName: tenant name cannot be empty")
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(c.host+pathTenantByName, name), nil)
	if err != nil {
		return errors.Wrap(err, "DeleteTenantByName: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return errors.Wrap(err, "DeleteTenantByName: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response code is not 200, got %d", res.StatusCode)
	}

	return nil
}
