package compute

import (
	"bytes"
	"encoding/json"
	"io"

	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func (c *Client) Exec(runnable *Runnable, body io.Reader) ([]byte, *http.Response, error) {
	req, err := c.execRequestBuilder(http.MethodPost, runnable.VersionPath(), body)

	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to Client.Exec")
	}

	res, err := c.do(req)
	if err != nil {
		return nil, res, errors.Wrap(err, "failed to Client.Exec")
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res, errors.Wrap(err, "failed to Client.Exec")
	}

	if res.StatusCode != http.StatusOK {
		errRes := ExecErrorResponse{}
		err = json.Unmarshal(result, &errRes)
		if err != nil {
			return nil, res, err
		}
		newErr := errors.Errorf("[%d]: %s", errRes.Code, errRes.Message)

		return nil, res, errors.Wrap(newErr, "failed to Client.Exec")
	}

	return result, res, nil
}

func (c *Client) ExecString(runnable *Runnable, body string) ([]byte, *http.Response, error) {
	buf := bytes.NewBufferString(body)
	return c.Exec(runnable, buf)
}
