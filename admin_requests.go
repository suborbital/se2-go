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

	runnable.token = token.Token

	return res, nil
}
