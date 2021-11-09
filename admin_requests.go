package compute

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

// EditorToken aquires and sets the editor token for the provided Runnable struct.
// See: Runnable.Token()
func (c *Client) EditorToken(runnable *Runnable) (*http.Response, error) {
	// GET /api/v1/token/{environment}.{customerID}/{namespace}/{fnName}
	req, err := c.adminRequestBuilder(http.MethodGet, path.Join("/api/v1/token", runnable.Path()), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return res, err
	}

	token := TokenResponse{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)

	if err != nil {
		return res, errors.Wrapf(err, "failed to getEditorTokenFor Runnable: [%s]", runnable.Path())
	}

	if token.Token == "" {
		return res, errors.Wrapf(errors.New("TokenReponse.Token was empty"),
			"failed to getEditorTokenFor Runnable: [%s]", runnable.Path())
	}

	runnable.editorToken = token.Token

	return res, nil
}

func (c *Client) UserFunctions(customerID string, namespace string) ([]*Runnable, *http.Response, error) {
	// GET /api/v1/functions/{customerID}/{namespace}

	req, err := c.adminRequestBuilder(http.MethodGet, path.Join("/api/v1/functions", customerID, namespace), nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, res, err
	}

	userFuncs := UserFunctionsResponse{
		Functions: []RunnableResponse{},
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&userFuncs)
	if err != nil {
		return nil, res, err
	}

	runnables := make([]*Runnable, len(userFuncs.Functions))
	for i, runnableRes := range userFuncs.Functions {
		runnables[i] = runnableRes.ToRunnable()
	}

	return runnables, res, nil
}

func (c *Client) FunctionExecResults(runnable *Runnable) (*ExecResultsResponse, *http.Response, error) {
	// GET /api/v1/results/com.awesomeco.vqeiupqvp98ph2e4nvrqw98/default/create-report/v0.0.1
	req, err := c.adminRequestBuilder(http.MethodGet, path.Join(
		"/api/v1/results", runnable.VersionPath()), nil)

	if err != nil {
		return nil, nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, res, err
	}

	execResults := &ExecResultsResponse{
		Results: []ExecResult{},
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(execResults)
	if err != nil {
		return nil, res, err
	}

	return execResults, res, nil
}
