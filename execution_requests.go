package se2

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

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
