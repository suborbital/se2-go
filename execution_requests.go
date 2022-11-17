package se2

import (
	"bytes"
	"encoding/json"
	"io"
	"path"

	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type ExecResponse struct {
	Response []byte
	UUID     string
}

// execModule remotely executes the provided module using the body as input. See also: ExecString()
func (c *Client) execModule(endpoint string, body io.Reader) ([]byte, string, error) {
	req, err := c.execRequestBuilder(http.MethodPost, endpoint, body)
	req.Header.Set("Authorization", "Bearer "+c.envToken)

	if err != nil {
		return nil, "", err
	}

	res, err := c.do(req)
	if err != nil && res == nil {
		return nil, "", err
	}
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	if res.StatusCode != http.StatusOK {
		errRes := ExecError{}
		err = json.Unmarshal(result, &errRes)
		if err != nil {
			return nil, "", err
		}
		newErr := errors.Errorf("[%d]: %s", errRes.Code, errRes.Message)

		return nil, "", newErr
	}

	return result, res.Header.Get("x-suborbital-requestid"), nil
}

// Exec remotely executes the provided module using the body as input. See also: ExecString()
func (c *Client) Exec(module *Module, body io.Reader) ([]byte, string, error) {
	if module == nil {
		return nil, "", errors.New("Module cannot be nil")
	}

	return c.execModule(path.Join("/name", module.URI()), body)
}

// ExecString sets up a buffer with the provided string and calls Exec
func (c *Client) ExecString(module *Module, body string) ([]byte, string, error) {
	buf := bytes.NewBufferString(body)
	return c.Exec(module, buf)
}

// ExecRef remotely executes the provided module using the body as input, by module reference. See also: ExecRefString()
func (c *Client) ExecRef(ref string, body io.Reader) ([]byte, string, error) {
	return c.execModule(path.Join("/ref", ref), body)
}

// ExecRefString sets up a buffer with the provided string and calls ExecRef
func (c *Client) ExecRefString(ref string, body string) ([]byte, string, error) {
	buf := bytes.NewBufferString(body)
	return c.ExecRef(ref, buf)
}
