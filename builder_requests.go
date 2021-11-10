package compute

import (
	"encoding/json"
	"net/http"
	"path"
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
