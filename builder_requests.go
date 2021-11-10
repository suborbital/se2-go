package compute

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"github.com/suborbital/atmo/directive"
)

func (c *Client) BuilderHealth() (bool, error) {
	req, err := c.builderRequestBuilder(http.MethodGet, "/api/v1/health", nil)
	if err != nil {
		return false, err
	}

	_, err = c.do(req)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) BuilderFeatures() (*FeaturesResponse, error) {
	req, err := c.builderRequestBuilder(http.MethodGet, "/api/v1/features", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	features := &FeaturesResponse{
		Features: []string{},
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(features)
	if err != nil {
		return nil, err
	}

	return features, nil
}

func (c *Client) BuilderTemplateV1(lang, namespace string) (*EditorStateResponse, error) {
	req, err := c.builderRequestBuilder(http.MethodGet,
		path.Join("/api/v1/template", lang), nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("namespace", namespace)
	req.URL.RawQuery = q.Encode()

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	editorState := &EditorStateResponse{Tests: []TestPayload{}}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(editorState)
	if err != nil {
		return nil, err
	}

	return editorState, nil
}

func (c *Client) BuilderTemplateV2(lang, namespace, fnName string) (*EditorStateResponse, error) {
	req, err := c.builderRequestBuilder(http.MethodGet,
		path.Join("/api/v2/template", lang, fnName), nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("namespace", namespace)
	req.URL.RawQuery = q.Encode()

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	editorState := &EditorStateResponse{Tests: []TestPayload{}}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(editorState)
	if err != nil {
		return nil, err
	}

	return editorState, nil
}

func (c *Client) BuildFunction(runnable *directive.Runnable, functionBody io.Reader) (*BuildResult, error) {
	if runnable == nil {
		return nil, errors.New("Runnable cannot be nil")
	}

	if functionBody == nil {
		return nil, errors.New("functionBody cannot be nil")
	}

	// TODO: cache somewhere in Client?
	token, err := c.EditorToken(runnable)
	if err != nil {
		return nil, errors.Wrap(err, "failed to EditorToken")
	}

	p, _ := path.Split(runnable.FQFNURI) // removes version from end of URI

	req, err := c.builderRequestBuilder(http.MethodPost,
		path.Join("/api/v1/build", runnable.Lang, p), functionBody)
	req.Header.Add("Authorization", "Bearer "+token)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	buildResult := &BuildResult{}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(buildResult)
	if err != nil {
		return nil, err
	}

	return buildResult, nil
}

func (c *Client) BuildFunctionString(runnable *directive.Runnable, functionString string) (*BuildResult, error) {
	buf := bytes.NewBufferString(functionString)
	return c.BuildFunction(runnable, buf)
}

func (c *Client) GetDraft(runnable *directive.Runnable) (*EditorStateResponse, error) {
	if runnable == nil {
		return nil, errors.New("Runnable cannot be nil")
	}

	token, err := c.EditorToken(runnable)
	if err != nil {
		return nil, errors.Wrap(err, "failed to EditorToken")
	}

	p, _ := path.Split(runnable.FQFNURI) // removes version from end of URI

	req, err := c.builderRequestBuilder(http.MethodGet,
		path.Join("/api/v1/draft", p), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	editorState := &EditorStateResponse{Tests: []TestPayload{}}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(editorState)
	if err != nil {
		return nil, err
	}

	return editorState, nil
}

func (c *Client) PromoteDraft(runnable *directive.Runnable) (*PromoteDraftResponse, error) {
	if runnable == nil {
		return nil, errors.New("Runnable cannot be nil")
	}

	token, err := c.EditorToken(runnable)
	if err != nil {
		return nil, errors.Wrap(err, "failed to EditorToken")
	}
	p, _ := path.Split(runnable.FQFNURI) // removes version from end of URI

	req, err := c.builderRequestBuilder(http.MethodPost,
		path.Join("/api/v1/draft", p, "promote"), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	promoteResponse := &PromoteDraftResponse{}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(promoteResponse)
	if err != nil {
		return nil, err
	}

	return promoteResponse, nil
}
