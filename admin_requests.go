package compute

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"github.com/suborbital/atmo/directive"
)

// EditorToken gets an editor token for the provided Runnable.
func (c *Client) EditorToken(runnable *directive.Runnable) (string, error) {
	// GET /api/v1/token/{environment}.{customerID}/{namespace}/{fnName}
	p, _ := path.Split(runnable.FQFNURI)
	req, err := c.adminRequestBuilder(http.MethodGet, path.Join("/api/v1/token", p), nil)
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
		return "", errors.Wrapf(err, "failed to getEditorTokenFor Runnable: [%s]", p)
	}

	if token.Token == "" {
		return "", errors.Wrapf(errors.New("TokenReponse.Token was empty"),
			"failed to getEditorTokenFor Runnable: [%s]", runnable.FQFN)
	}

	return token.Token, nil
}

func (c *Client) UserFunctions(customerID string, namespace string) ([]*directive.Runnable, error) {
	// GET /api/v1/functions/{customerID}/{namespace}

	req, err := c.adminRequestBuilder(http.MethodGet, path.Join("/api/v1/functions", customerID, namespace), nil)
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

func (c *Client) FunctionExecResults(runnable *directive.Runnable) (*ExecResultsResponse, error) {
	// GET /api/v1/results/com.awesomeco.vqeiupqvp98ph2e4nvrqw98/default/create-report/v0.0.1
	req, err := c.adminRequestBuilder(http.MethodGet, path.Join(
		"/api/v1/results", runnable.FQFNURI), nil)

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
