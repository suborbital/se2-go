package compute

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"github.com/suborbital/atmo/directive"
)

// EditorToken gets an editor token for the provided Runnable. Note: this library
// manages editor tokens for you, so you most likely do not need to use this function.
func (c *Client) EditorToken(runnable *directive.Runnable) (string, error) {
	if runnable == nil {
		return "", errors.New("Runnable cannot be nil")
	}

	p, _ := path.Split(runnable.FQFNURI) // removes version from end of URI
	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v1/token", p), nil)

	if err != nil {
		return "", err
	}

	res, err := c.do(req)
	if err != nil {
		return "", err
	}

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

// UserFunctions gets a list of the deployed runnables for the given customer and namespace.
func (c *Client) UserFunctions(customerID string, namespace string) ([]*directive.Runnable, error) {
	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v1/functions", customerID, namespace), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	userFuncs := UserFunctionsResponse{
		Functions: []*directive.Runnable{},
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&userFuncs)
	if err != nil {
		return nil, err
	}

	return userFuncs.Functions, nil
}

// FunctionExecResults returns the 5 most recent successful execution results for the provided runnable.
func (c *Client) FunctionExecResults(runnable *directive.Runnable) (*ExecResultsResponse, error) {
	if runnable == nil {
		return nil, errors.New("Runnable cannot be nil")
	}

	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v1/results", runnable.FQFNURI), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	execResults := &ExecResultsResponse{
		Results: []ExecResult{},
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(execResults)
	if err != nil {
		return nil, err
	}

	return execResults, nil
}

// FunctionExecResults returns the 5 most recent execution errors for the provided runnable.
func (c *Client) FunctionExecErrors(runnable *directive.Runnable) (*ExecErrorResponse, error) {
	if runnable == nil {
		return nil, errors.New("Runnable cannot be nil")
	}

	req, err := c.adminRequestBuilder(http.MethodGet,
		path.Join("/api/v1/errors", runnable.FQFNURI), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	execResults := &ExecErrorResponse{
		Errors: []ExecError{},
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(execResults)
	if err != nil {
		return nil, err
	}

	return execResults, nil
}
