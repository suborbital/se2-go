package se2

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"github.com/suborbital/systemspec/tenant"
)

// EditorToken gets an editor token for the provided Plugin. Note: this library
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
func (c *Client) UserPlugins(identifier string, namespace string) ([]*tenant.Module, error) {
	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v2/functions", identifier, namespace), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	userPlugins := UserPluginsResponse{
		Plugins: []*tenant.Module{},
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&userPlugins)
	if err != nil {
		return nil, err
	}

	return userPlugins.Plugins, nil
}

// ResultsMetadata returns metadata for the 5 most recent execution results for the provided plugin.
func (c *Client) ResultsMetadata(plugin *Plugin) ([]ExecMetadata, error) {
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

// ResultMetadata returns metadata for the provided plugin execution.
func (c *Client) ResultMetadata(uuid string) (*ExecMetadata, error) {
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

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
