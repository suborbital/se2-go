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

func (c *Client) Exec(runnable *directive.Runnable, body io.Reader) ([]byte, error) {
	if runnable == nil {
		return nil, errors.New("Runnable cannot be nil")
	}

	req, err := c.execRequestBuilder(http.MethodPost, runnable.FQFNURI, body)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil && res == nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		errRes := ExecError{}
		err = json.Unmarshal(result, &errRes)
		if err != nil {
			return nil, err
		}
		newErr := errors.Errorf("[%d]: %s", errRes.Code, errRes.Message)

		return nil, newErr
	}

	return result, nil
}

func (c *Client) ExecString(runnable *directive.Runnable, body string) ([]byte, error) {
	buf := bytes.NewBufferString(body)
	return c.Exec(runnable, buf)
}
