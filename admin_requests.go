package se2

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"github.com/suborbital/systemspec/tenant"
)

// EditorToken gets an editor token for the provided Module. Note: this library
// manages editor tokens for you, so you most likely do not need to use this function.
func (c *Client) EditorToken(module *Module) (string, error) {
	if module == nil {
		return "", errors.New("Module cannot be nil")
	}

	req, err := c.adminRequestBuilder(http.MethodGet, path.Join("/api/v1/token", module.URI()), nil)

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

// UserFunctions gets a list of the deployed modules for the given identifier and namespace.
func (c *Client) UserFunctions(identifier string, namespace string) ([]*tenant.Module, error) {
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

	userFuncs := UserFunctionsResponse{
		Functions: []*tenant.Module{},
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&userFuncs)
	if err != nil {
		return nil, err
	}

	return userFuncs.Functions, nil
}

// FunctionResultsMetadata returns metadata for the 5 most recent execution results for the provided module.
func (c *Client) FunctionResultsMetadata(module *Module) ([]ExecMetadata, error) {
	if module == nil {
		return nil, errors.New("Module cannot be nil")
	}

	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v2/results/by-fqmn", module.URI()), nil)

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

// FunctionResultMetadata returns metadata for the provided module execution.
func (c *Client) FunctionResultMetadata(uuid string) (*ExecMetadata, error) {
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

// FunctionResult returns the result of the provided module execution.
func (c *Client) FunctionResult(uuid string) ([]byte, error) {
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
