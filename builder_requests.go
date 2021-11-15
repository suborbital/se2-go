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

// BuilderHealth is used to check that the builder is healthy and responding to requests.
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

// BuilderTemplate gets the function template for the provided Runnable. The Runnable must have
// the .Lang, .Name, and .Namespace fields set.
func (c *Client) BuilderTemplate(runnable *directive.Runnable) (*EditorStateResponse, error) {
	if runnable == nil {
		return nil, errors.New("Runnable cannot be nil")
	}

	if runnable.Lang == "" || runnable.Name == "" || runnable.Namespace == "" {
		return nil, errors.New("Runnable.Lang, .Name, and .Namespace must be set")
	}

	req, err := c.builderRequestBuilder(http.MethodGet,
		path.Join("/api/v2/template", runnable.Lang, runnable.Name), nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("namespace", runnable.Namespace)
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

// BuildFunction triggers a remote build for the given Runnable and function body. See also: Client.BuildFunctionString()
//
// Example
//
// This function is useful for reading from a filesystem or from an http.Response.Body
//	runnable := compute.NewRunnable("com.suborbital", "acmeco", "default", "hello", "rust")
// 	file, _ := os.Open("hello.rs")
//	result, err := client.BuildFunction(runnable, file)
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

// BuildFunctionString triggers a remote build for the given Runnable and function string. See also: Client.BuildFunction()
func (c *Client) BuildFunctionString(runnable *directive.Runnable, functionString string) (*BuildResult, error) {
	buf := bytes.NewBufferString(functionString)
	return c.BuildFunction(runnable, buf)
}

// GetDraft gets the most recently build source code for the provided Runnable. Must have the .FQFNURI field set.
func (c *Client) GetDraft(runnable *directive.Runnable) (*EditorStateResponse, error) {
	if runnable == nil {
		return nil, errors.New("Runnable cannot be nil")
	}

	if runnable.FQFNURI == "" {
		return nil, errors.New("Runnable.FQFNURI must be set")
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

// PromoteDraft takes the most recent build of the provided runnable and deploys it so it can be
// run. The .Version field of the provided runnable is modified in place if the promotion is
// successful.
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

	runnable.Version = promoteResponse.Version

	return promoteResponse, nil
}
