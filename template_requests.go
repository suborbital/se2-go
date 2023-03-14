package se2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	pathTemplate       = "/template/v1"
	pathTemplateByName = pathTemplate + "/%s"
	pathTemplateImport = pathTemplate + "/import"
)

// ListTemplatesResponse is used to marshal returned json from the SE2 backend into a struct.
type ListTemplatesResponse struct {
	Templates []Template `json:"templates"`
}

// Template holds information about a singular template.
type Template struct {
	Name    string `json:"name"`
	Lang    string `json:"lang"`
	Main    string `json:"main,omitempty"`
	Version string `json:"api_version"`
}

// ListTemplates will return a ListTemplatesResponse which contains a slice of Template that are available to the
// environment specified by the API key of the client.
func (c *Client) ListTemplates(ctx context.Context) (ListTemplatesResponse, error) {
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

	var t ListTemplatesResponse
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return ListTemplatesResponse{}, errors.Wrap(err, "ListTemplates: dec.Decode")
	}

	return t, nil
}

// GetTemplate takes a name and will return information about a template by that name, or an error if no templates are
// found.
func (c *Client) GetTemplate(ctx context.Context, name string) (Template, error) {
	if name == emptyString {
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

	var t Template
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return Template{}, errors.Wrap(err, "GetTemplate: dec.Decode")
	}

	return t, nil
}

// importRequest is an internal shape of the data we need to send to the SE2 backend.
type importRequest struct {
	Source string       `json:"source"`
	Params importParams `json:"params"`
}

// importParams is an internal shape to help describe data to send to the import endpoint on the API.
type importParams struct {
	Repo string `json:"repo"`
	Ref  string `json:"ref"`
	Path string `json:"path"`
}

// ImportTemplatesFromGitHub takes a repo, ref, and a path as arguments. It will call the appropriate endpoint on the
// configured SE2 backend. The values are ultimately used to download the source code of a given repository at a given
// reference point - either commit sha or tag reference.
//
// https://github.com/{repo}/archive/{ref}.tar.gz is the pattern.
// ref can take shape in one of the following ways:
// - 755721878465402d2f574635f8b411172f9f8482 (example) if it's a commit sha. At least 6 characters usually works well
// - refs/tags/v0.3.2 if it's a reference to a tag
//
// The repository needs to be publicly accessible; private repositories are not supported. Right now only GitHub is the
// only available provider we can pull source code from.
func (c *Client) ImportTemplatesFromGitHub(ctx context.Context, repo, ref, path string) error {
	if repo == emptyString {
		return errors.New("repo cannot be blank")
	}

	if ref == emptyString {
		return errors.New("ref cannot be blank")
	}

	if path == emptyString {
		return errors.New("path cannot be blank. If files are on the root of the repository, use '.'")
	}

	var requestBody bytes.Buffer
	err := json.NewEncoder(&requestBody).Encode(importRequest{
		Source: "git",
		Params: importParams{
			Repo: repo,
			Ref:  ref,
			Path: path,
		},
	})
	if err != nil {
		return errors.Wrap(err, "ImportTemplatesFromGit: json.NewEncoder.Encode")
	}

	req, err := http.NewRequest(http.MethodPost, c.host+pathTemplateImport, &requestBody)
	if err != nil {
		return errors.Wrap(err, "ImportTemplatesFromGit: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return errors.Wrap(err, "ImportTemplatesFromGit: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("ImportTemplatesFromGit: expected response to be %d, got %d instead", http.StatusOK, res.StatusCode)
	}

	return nil
}
