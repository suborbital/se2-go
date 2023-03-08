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
	pathTemplate       = "/template/v1"
	pathTemplateByName = pathTemplate + "/%s"
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
	req, err := http.NewRequest(http.MethodGet, c.host+pathTemplate, nil)
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
		return ListTemplatesResponse{}, fmt.Errorf("ListTemplates: unexpected status code. Wanted %d, got %d", http.StatusOK, res.StatusCode)
	}

	b, _ := io.ReadAll(res.Body)
	fmt.Printf("%s\n\n", string(b))

	var t ListTemplatesResponse

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return ListTemplatesResponse{}, errors.Wrap(err, "ListTemplates: dec.Decode")
	}

	return t, nil
}

func (c *Client2) GetTemplate(ctx context.Context, name string) (Template, error) {
	if name == "" {
		return Template{}, errors.New("name cannot be blank")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(c.host+pathTemplateByName, name), nil)
	if err != nil {
		return Template{}, errors.Wrap(err, "GetTemplate: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return Template{}, errors.Wrap(err, "GetTemplate: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	b, _ := io.ReadAll(res.Body)

	fmt.Printf("gettemplate body:\n%s\n\n", string(b))

	var t Template

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return Template{}, errors.Wrap(err, "GetTemplate: dec.Decode")
	}

	return t, nil
}
