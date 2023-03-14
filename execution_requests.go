package se2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

type ExecResponse struct {
	Response []byte
}

const (
	pathExec = "/name/%s/%s/%s"
)

// Exec takes a context, a byte slice payload, an ident, namespace, and plugin triad to identify the plugin to run with
// the payload as input. It returns a byte slice as output, and an error if something went wrong.
func (c *Client2) Exec(ctx context.Context, payload []byte, ident, namespace, plugin string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(c.execHost+pathExec, ident, namespace, plugin), bytes.NewReader(payload))
	if err != nil {
		return nil, errors.Wrap(err, "Exec: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "Exec: c.builderDo")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Exec: io.ReadAll(res.Body)")
	}

	return b, nil
}

// execPlugin remotely executes the provided plugin using the body as input. See also: ExecString()
func (c *Client) execPlugin(endpoint string, body io.Reader) ([]byte, string, error) {
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

	result, err := io.ReadAll(res.Body)
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

// Exec remotely executes the provided plugin using the body as input. See also: ExecString()
func (c *Client) Exec(plugin *Plugin, body io.Reader) ([]byte, string, error) {
	if plugin == nil {
		return nil, "", errors.New("Plugin cannot be nil")
	}

	return c.execPlugin(path.Join("/name", plugin.URI()), body)
}

// ExecString sets up a buffer with the provided string and calls Exec
func (c *Client) ExecString(plugin *Plugin, body string) ([]byte, string, error) {
	buf := bytes.NewBufferString(body)
	return c.Exec(plugin, buf)
}

// ExecRef remotely executes the provided plugin using the body as input, by plugin reference. See also: ExecRefString()
func (c *Client) ExecRef(ref string, body io.Reader) ([]byte, string, error) {
	return c.execPlugin(path.Join("/ref", ref), body)
}

// ExecRefString sets up a buffer with the provided string and calls ExecRef
func (c *Client) ExecRefString(ref string, body string) ([]byte, string, error) {
	buf := bytes.NewBufferString(body)
	return c.ExecRef(ref, buf)
}
