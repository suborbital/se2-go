package se2

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"path"

	"github.com/pkg/errors"
	systemspec "github.com/suborbital/systemspec/tenant"
)

// EditorToken gets an editor token for the provided plugin. Note: this library
// manages editor tokens for you, so you most likely do not need to use this function.
func (c *Client) EditorToken(plugin *Plugin) (string, error) {
	if plugin == nil {
		return "", errors.New("Plugin cannot be nil")
	}

	req, err := c.adminRequestBuilder(http.MethodGet, path.Join("/api/v1/token", plugin.URI()), nil)

	if err != nil {
		return "", err
	}

	res, err := c.do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	token := TokenResponse{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)

	if err != nil {
		return "", err
	}

	if token.Token == "" {
		return "", err
	}

	return token.Token, nil
}

// UserPlugins gets a list of the deployed plugins for the given identifier and namespace.
func (c *Client) UserPlugins(environment string, tenant string, namespace string) ([]*Plugin, error) {
	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v2/functions", fmt.Sprintf("%s.%s", environment, tenant), namespace), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	userPlugins := UserPluginsResponse{
		Plugins: []*systemspec.Module{},
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&userPlugins)
	if err != nil {
		return nil, err
	}

	plugins := make([]*Plugin, len(userPlugins.Plugins))
	for i, module := range userPlugins.Plugins {
		plugins[i].Environment = environment
		plugins[i].Tenant = tenant
		plugins[i].Namespace = namespace
		plugins[i].Name = module.Name
	}

	return plugins, nil
}

// ExecutionResultsMetadata returns metadata for the 5 most recent execution results for the provided plugin.
func (c *Client) ExecutionResultsMetadata(plugin *Plugin) ([]ExecMetadata, error) {
	if plugin == nil {
		return nil, errors.New("Plugin cannot be nil")
	}

	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v2/results/by-fqmn", plugin.URI()), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var execResults []ExecMetadata

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&execResults)
	if err != nil {
		return nil, err
	}

	return execResults, nil
}

// ExecutionResultMetadata returns metadata for the provided plugin execution.
func (c *Client) ExecutionResultMetadata(uuid string) (*ExecMetadata, error) {
	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v2/results/by-uuid", uuid), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var execResult ExecMetadata

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&execResult)
	if err != nil {
		return nil, err
	}

	return &execResult, nil
}

// ExecutionResult returns the result of the provided plugin execution.
func (c *Client) ExecutionResult(uuid string) ([]byte, error) {
	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v2/result", uuid), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
