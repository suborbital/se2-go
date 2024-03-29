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
	pathBuilderFeatures = "/api/v1/features"
	pathBuilderPrefix   = "/builder/v1"
	pathDraft           = pathBuilderPrefix + "/draft"
	pathBuild           = pathDraft + "/build"
	pathTest            = pathDraft + "/test"
	pathPromote         = pathDraft + "/deploy"
)

// BuildPluginResponse captures the json body into a struct for a  build attempt.
type BuildPluginResponse struct {
	Succeeded bool   `json:"succeeded"`
	OutputLog string `json:"outputLog"`
}

// BuildPlugin will attempt to build a plugin supplied by the raw byte slice in the context of the current session. The
// language is set by the template, which you can control by calling the CreatePluginDraft method with the template
// name.
func (c *Client) BuildPlugin(ctx context.Context, pluginCode []byte, token CreateSessionResponse) (BuildPluginResponse, error) {
	if len(pluginCode) == zeroLength {
		return BuildPluginResponse{}, errors.New("client.BuildPlugin: can not build empty code")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+pathBuild, bytes.NewReader(pluginCode))
	if err != nil {
		return BuildPluginResponse{}, errors.Wrap(err, "client.BuildPlugin: http.NewRequest")
	}

	res, err := c.sessionDo(req, token)
	if err != nil {
		return BuildPluginResponse{}, errors.Wrap(err, "client.BuildPlugin: c.sessionDo")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusCreated {
		return BuildPluginResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.BuildPlugin", http.StatusCreated, res.StatusCode)
	}

	var t BuildPluginResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return BuildPluginResponse{}, errors.Wrap(err, "client.BuildPlugin: dec.Decode")
	}

	return t, nil
}

// BuilderFeaturesResponse captures the json response from the features endpoint.
type BuilderFeaturesResponse struct {
	Features  []string `json:"features"`
	Languages []Languages
}

// Languages captures the json representation of an individual supported language.
type Languages struct {
	ID         string `json:"identifier"`
	ShortName  string `json:"short"`
	PrettyName string `json:"pretty"`
}

// GetBuilderFeatures will return the features that the builder can provide.
func (c *Client) GetBuilderFeatures(ctx context.Context) (BuilderFeaturesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.host+pathBuilderFeatures, nil)
	if err != nil {
		return BuilderFeaturesResponse{}, errors.Wrap(err, "client.GetBuilderFeatures: http.NewRequest")
	}

	res, err := c.do(req)
	if err != nil {
		return BuilderFeaturesResponse{}, errors.Wrap(err, "client.GetBuilderFeatures: c.do")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return BuilderFeaturesResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.GetBuilderFeatures", http.StatusOK, res.StatusCode)
	}

	// Marshal response body into what we need to give back.
	var t BuilderFeaturesResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return BuilderFeaturesResponse{}, errors.Wrap(err, "client.GetBuilderFeatures: dec.Decode")
	}

	return t, nil
}

// runError captures the structure that the Error in the test draft response can take.
type runError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// TestPluginDraftResponse is the response of the test call with the given input data.
type TestPluginDraftResponse struct {
	Result string   `json:"result"`
	Error  runError `json:"error"`
}

// TestPluginDraft will send the testData byte slice to the plugin that's currently in the draft as input, and return
// the response that came back from the plugin.
func (c *Client) TestPluginDraft(ctx context.Context, testData []byte, token CreateSessionResponse) (TestPluginDraftResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+pathTest, bytes.NewReader(testData))
	if err != nil {
		return TestPluginDraftResponse{}, errors.Wrap(err, "client.TestPluginDraft: http.NewRequest")
	}

	res, err := c.sessionDo(req, token)
	if err != nil {
		return TestPluginDraftResponse{}, errors.Wrap(err, "client.TestPluginDraft: c.sessionDo")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return TestPluginDraftResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.TestPluginDraft", http.StatusOK, res.StatusCode)
	}

	var t TestPluginDraftResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return TestPluginDraftResponse{}, errors.Wrap(err, "client.TestPluginDraft: dec.Decode")
	}

	return t, nil
}

// GetPluginDraft returns the currently set plugin draft for the given session token. To change the draft or the
// language you can use the CreatePluginDraft method instead with the name of a template.
func (c *Client) GetPluginDraft(ctx context.Context, token CreateSessionResponse) (DraftResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.host+pathDraft, nil)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "client.GetPluginDraft: http.NewRequest")
	}

	res, err := c.sessionDo(req, token)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "client.GetPluginDraft: c.sessionDo")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return DraftResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.GetPluginDraft", http.StatusOK, res.StatusCode)
	}

	var t DraftResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "client.GetPluginDraft: dec.Decode")
	}

	return t, nil
}

// DraftResponse is a struct the captures the response from the CreatePluginDraft and GetDraft endpoints.
type DraftResponse struct {
	Lang     string `json:"lang"`
	Contents string `json:"contents"`
}

// createDraftRequest is a helper struct to encode an incoming template name into a correct json structure we can send
// to the API. Users of this client library should not need to interact with this struct directly.
type createDraftRequest struct {
	Template string `json:"template"`
}

// CreatePluginDraft takes in a template name and a session token, and will set the current session to use the named
// template for building and executing.
//
// To see available templates, use the ListTemplates method.
func (c *Client) CreatePluginDraft(ctx context.Context, templateName string, token CreateSessionResponse) (DraftResponse, error) {
	if templateName == emptyString {
		return DraftResponse{}, errors.New("client.CreatePluginDraft: template name cannot be blank")
	}

	var b bytes.Buffer

	r := createDraftRequest{Template: templateName}

	err := json.NewEncoder(&b).Encode(r)
	if err != nil {
		return DraftResponse{}, errors.Wrapf(err, "client.CreatePluginDraft: json.NewEncoder.Encode(createDraftRequest with template name '%s'", templateName)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+pathDraft, &b)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "client.CreatePluginDraft: http.NewRequest")
	}

	res, err := c.sessionDo(req, token)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "client.CreatePluginDraft: c.sessionDo")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return DraftResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.CreatePluginDraft", http.StatusOK, res.StatusCode)
	}

	var t DraftResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "client.CreatePluginDraft: dec.Decode")
	}

	return t, nil
}

// PromotePluginDraftResponse captures the json returned by a successful call to the promote endpoint.
type PromotePluginDraftResponse struct {
	Ref string `json:"ref"`
}

// PromotePluginDraft promotes the current version of the draft to the live version of the plugin.
func (c *Client) PromotePluginDraft(ctx context.Context, token CreateSessionResponse) (PromotePluginDraftResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+pathPromote, nil)
	if err != nil {
		return PromotePluginDraftResponse{}, errors.Wrap(err, "client.PromotePluginDraft: http.NewRequest")
	}

	res, err := c.sessionDo(req, token)
	if err != nil {
		return PromotePluginDraftResponse{}, errors.Wrap(err, "client.PromotePluginDraft: c.sessionDo")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return PromotePluginDraftResponse{}, fmt.Errorf(httpResponseCodeErrorFormat, "client.PromotePluginDraft", http.StatusOK, res.StatusCode)
	}

	var t PromotePluginDraftResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&t)
	if err != nil {
		return PromotePluginDraftResponse{}, errors.Wrap(err, "client.PromotePluginDraft: dec.Decode")
	}

	return t, nil
}
