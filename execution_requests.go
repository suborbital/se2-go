package compute

import (
	"bytes"
	"encoding/json"
	"io"

	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/suborbital/atmo/directive"
)

type ExecResponse struct {
	Response []byte
	UUID     string
}

// Exec remotely executes the provided runnable using the body as input. See also: ExecString()
func (c *Client) Exec(runnable *directive.Runnable, body io.Reader) ([]byte, string, error) {
	if runnable == nil {
		return nil, "", errors.New("Runnable cannot be nil")
	}

	req, err := c.execRequestBuilder(http.MethodPost, runnable.FQFNURI, body)
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

	return result, res.Header.Get("x-atmo-requestid"), nil
}

// ExecString sets up a buffer with the provided string and calls Exec
func (c *Client) ExecString(runnable *directive.Runnable, body string) ([]byte, string, error) {
	buf := bytes.NewBufferString(body)
	return c.Exec(runnable, buf)
}
