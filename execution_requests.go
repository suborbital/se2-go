package se2

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	pathExec = "/name/%s/%s/%s"
)

// Exec takes a context, a byte slice payload, an ident, namespace, and plugin triad to identify the plugin to run with
// the payload as input. It returns a byte slice as output, and an error if something went wrong.
func (c *Client) Exec(ctx context.Context, payload []byte, ident, namespace, plugin string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(c.execHost+pathExec, ident, namespace, plugin), bytes.NewReader(payload))
	if err != nil {
		return nil, errors.Wrap(err, "client.Exec: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Exec: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(httpResponseCodeErrorFormat, "client.Exec", http.StatusOK, res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "client.Exec: io.ReadAll(res.Body)")
	}

	return b, nil
}
