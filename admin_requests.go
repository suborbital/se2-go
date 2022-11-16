package compute

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"github.com/suborbital/systemspec/tenant"
)

// EditorToken gets an editor token for the provided Runnable. Note: this library
// manages editor tokens for you, so you most likely do not need to use this function.
func (c *Client) EditorToken(runnable *tenant.Module) (string, error) {
	if runnable == nil {
		return "", errors.New("Runnable cannot be nil")
	}

	p, _ := path.Split(runnable.URI) // removes version from end of URI
	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v1/token", p), nil)

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

// UserFunctions gets a list of the deployed runnables for the given identifier and namespace.
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

// FunctionResultsMetadata returns metadata for the 5 most recent execution results for the provided runnable.
func (c *Client) FunctionResultsMetadata(runnable *tenant.Module) ([]ExecMetadata, error) {
	if runnable == nil {
		return nil, errors.New("Runnable cannot be nil")
	}

	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v2/results/by-fqfn", runnable.URI), nil)

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

// FunctionResultMetadata returns metadata for the provided runnable execution.
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

// FunctionResult returns the result of the provided runnable execution.
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
