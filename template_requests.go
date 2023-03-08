package se2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	pathTemplates = "/template/v1"
)

type ListTemplatesResponse struct {
	Templates []Template `json:"templates"`
}

type Template struct {
	Name    string `json:"name"`
	Lang    string `json:"lang"`
	Main    string `json:"main,omitempty"`
	Version string `json:"api_version"`
}

func (c *Client2) ListTemplates(ctx context.Context) (ListTemplatesResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.host+pathTemplates, nil)
	if err != nil {
		return ListTemplatesResponse{}, errors.Wrap(err, "ListTemplates: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return ListTemplatesResponse{}, errors.Wrap(err, "ListTemplates: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return ListTemplatesResponse{}, fmt.Errorf("unexpected status code. Wanted %d, got %d", http.StatusOK, res.StatusCode)
	}

	b, _ := io.ReadAll(res.Body)
	fmt.Printf("%s\n\n", string(b))

	var t ListTemplatesResponse

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return ListTemplatesResponse{}, errors.Wrap(err, "GetBuilderFeatures: dec.Decode")
	}

	return t, nil
}
