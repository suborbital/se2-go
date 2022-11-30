package se2

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

// BuilderHealth is used to check that the builder is healthy and responding to requests.
func (c *Client) BuilderHealth() (bool, error) {
	req, err := c.builderRequestBuilder(http.MethodGet, "/api/v1/health", nil)
	if err != nil {
		return false, err
	}

	res, err := c.do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return true, nil
}

// BuilderFeatures lists the features present on the builder, such as testing capabilities.
func (c *Client) BuilderFeatures() (*FeaturesResponse, error) {
	req, err := c.builderRequestBuilder(http.MethodGet, "/api/v1/features", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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

// BuilderTemplate gets the function template for the provided Plugin and template name.
func (c *Client) BuilderTemplate(plugin *Plugin, template string) (*EditorStateResponse, error) {
	if plugin == nil {
		return nil, errors.New("Plugin cannot be nil")
	}

	req, err := c.builderRequestBuilder(http.MethodGet,
		path.Join("/api/v2/template", template, plugin.Name), nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("namespace", plugin.Namespace)
	req.URL.RawQuery = q.Encode()

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	editorState := &EditorStateResponse{Tests: []TestPayload{}}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(editorState)
	if err != nil {
		return nil, err
	}

	return editorState, nil
}

// BuildFunction triggers a remote build for the given Plugin and function body. See also: Client.BuildFunctionString()
//
// # Example
//
// This function is useful for reading from a filesystem or from an http.Response.Body
//
//	plugin := compute.NewPlugin("com.suborbital", "acmeco", "default", "hello", "rust")
//	file, _ := os.Open("hello.rs")
//	result, err := client.BuildFunction(plugin, file)
func (c *Client) BuildFunction(plugin *Plugin, template string, functionBody io.Reader) (*BuildResult, error) {
	if plugin == nil {
		return nil, errors.New("Plugin cannot be nil")
	}

	if functionBody == nil {
		return nil, errors.New("functionBody cannot be nil")
	}

	// TODO: cache somewhere in Client?
	token, err := c.EditorToken(plugin)
	if err != nil {
		return nil, errors.Wrap(err, "failed to EditorToken")
	}

	req, err := c.builderRequestBuilder(http.MethodPost,
		path.Join("/api/v1/build", template, plugin.URI()), functionBody)
	req.Header.Add("Authorization", "Bearer "+token)

	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buildResult := &BuildResult{}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(buildResult)
	if err != nil {
		return nil, err
	}

	return buildResult, nil
}

// BuildFunctionString triggers a remote build for the given Plugin and function string. See also: Client.BuildFunction()
func (c *Client) BuildFunctionString(plugin *Plugin, template, functionString string) (*BuildResult, error) {
	buf := bytes.NewBufferString(functionString)
	return c.BuildFunction(plugin, template, buf)
}

// GetDraft gets the most recently build source code for the provided Plugin. Must have the .FQFNURI field set.
func (c *Client) GetDraft(plugin *Plugin) (*EditorStateResponse, error) {
	if plugin == nil {
		return nil, errors.New("Plugin cannot be nil")
	}

	token, err := c.EditorToken(plugin)
	if err != nil {
		return nil, errors.Wrap(err, "failed to EditorToken")
	}

	req, err := c.builderRequestBuilder(http.MethodGet,
		path.Join("/api/v1/draft", plugin.URI()), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	editorState := &EditorStateResponse{Tests: []TestPayload{}}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(editorState)
	if err != nil {
		return nil, err
	}

	return editorState, nil
}

// PromoteDraft takes the most recent build of the provided plugin and deploys it so it can be
// run.
func (c *Client) PromoteDraft(plugin *Plugin) (*PromoteDraftResponse, error) {
	if plugin == nil {
		return nil, errors.New("Plugin cannot be nil")
	}

	token, err := c.EditorToken(plugin)
	if err != nil {
		return nil, errors.Wrap(err, "failed to EditorToken")
	}

	req, err := c.builderRequestBuilder(http.MethodPost,
		path.Join("/api/v1/draft", plugin.URI(), "promote"), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	promoteResponse := &PromoteDraftResponse{}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(promoteResponse)
	if err != nil {
		return nil, err
	}

	return promoteResponse, nil
}
