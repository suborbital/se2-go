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
	ID              string `json:"id"`
	Environment     string `json:"environment"`
	Name            string `json:"name"`
	Description     string `json:"description"`
}

// GetTenantByName returns the tenant by name.
func (c *Client) GetTenantByName(ctx context.Context, name string) (TenantResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.host+fmt.Sprintf(pathTenantByName, name), nil)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.GetTenantByName: http.NewRequest")
	}

	res, err := c.do(req)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.GetTenantByName: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return TenantResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.GetTenantByName", http.StatusOK, res.StatusCode)
	}

	var t TenantResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.GetTenantByName: dec.Decode")
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
func (c *Client) CreateTenant(ctx context.Context, name, description string) (TenantResponse, error) {
	if name == emptyString {
		return TenantResponse{}, errors.New("client.CreateTenant: tenant name cannot be empty")
	}

	var requestBody io.Reader

	if description != "" {
		m, err := json.Marshal(createTenantRequest{Description: description})
		if err != nil {
			return TenantResponse{}, errors.Wrap(err, "client.CreateTenant: json marshal create tenant request with description")
		}

		requestBody = bytes.NewReader(m)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+fmt.Sprintf(pathTenantByName, name), requestBody)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.CreateTenant: http.NewRequest for POST create tenant")
	}

	res, err := c.do(req)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.CreateTenant: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusCreated {
		return TenantResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.CreateTenant", http.StatusCreated, res.StatusCode)
	}

	var t TenantResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.CreateTenant: dec.Decode")
	}

	return t, nil
}

// ListTenantResponse is the unmarshaled response from the endpoint.
type ListTenantResponse struct {
	Tenants []TenantResponse
}

// ListTenants will list the tenants that the configured API key can access.
func (c *Client) ListTenants(ctx context.Context) (ListTenantResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.host+pathTenant, nil)
	if err != nil {
		return ListTenantResponse{}, errors.Wrap(err, "client.ListTenants: http.NewRequest")
	}

	res, err := c.do(req)
	if err != nil {
		return ListTenantResponse{}, errors.Wrap(err, "client.ListTenants: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return ListTenantResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.ListTenants", http.StatusOK, res.StatusCode)
	}

	var t ListTenantResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return ListTenantResponse{}, errors.Wrap(err, "client.ListTenants: dec.Decode")
	}

	return t, nil
}

// updateTenantRequest is a struct to help constrain the data and marshal the json based on it that ends up in a request
// body. It's internal only, users of the client do not need to know about the existence of this struct.
type updateTenantRequest struct {
	Description string `json:"description"`
}

// UpdateTenantByName updates the description of the tenant identified by its name. A tenant's name cannot be changed.
func (c *Client) UpdateTenantByName(ctx context.Context, name, description string) (TenantResponse, error) {
	if name == emptyString {
		return TenantResponse{}, errors.New("client.UpdateTenantByName: tenant name cannot be empty")
	}

	m, err := json.Marshal(updateTenantRequest{Description: description})
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.UpdateTenantByName: json marshal update tenant request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.host+fmt.Sprintf(pathTenantByName, name), bytes.NewReader(m))
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.UpdateTenantByName: http.NewRequest for POST create tenant")
	}

	res, err := c.do(req)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.UpdateTenantByName: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return TenantResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.UpdateTenantByName", http.StatusOK, res.StatusCode)
	}

	var t TenantResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return TenantResponse{}, errors.Wrap(err, "client.UpdateTenantByName: dec.Decode")
	}

	return t, nil
}

// DeleteTenantByName deletes the tenant identified by its name.
func (c *Client) DeleteTenantByName(ctx context.Context, name string) error {
	if name == emptyString {
		return errors.New("client.DeleteTenantByName: tenant name cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf(c.host+pathTenantByName, name), nil)
	if err != nil {
		return errors.Wrap(err, "client.DeleteTenantByName: http.NewRequest")
	}

	res, err := c.do(req)
	if err != nil {
		return errors.Wrap(err, "client.DeleteTenantByName: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(httpResponseCodeErrorFormat, "client.DeleteTenantByName", http.StatusOK, res.StatusCode)
	}

	return nil
}
