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
	pathBuilderPrefix      = "/builder/v1"
	pathBuilderFeatures    = pathBuilderPrefix + "/features"
	pathBuilderFeaturesOld = "/api/v1/features"
	pathDraft              = pathBuilderPrefix + "/draft"
	pathBuild              = pathDraft + "/build"
)

type BuildPluginRequest struct{}

type BuildPluginResponse struct {
	Succeeded bool   `json:"succeeded"`
	OutputLog string `json:"outputLog"`
}

func (c *Client2) BuildPlugin(ctx context.Context, pluginCode []byte, token CreateSessionResponse) (BuildPluginResponse, error) {
	if len(pluginCode) == 0 {
		return BuildPluginResponse{}, errors.New("can not build empty code")
	}

	req, err := http.NewRequest(http.MethodPost, c.host+pathBuild, bytes.NewReader(pluginCode))
	if err != nil {
		return BuildPluginResponse{}, errors.Wrap(err, "BuildPlugin: http.NewRequest")
	}

	res, err := c.builderDo(ctx, req, token)
	if err != nil {
		return BuildPluginResponse{}, errors.Wrap(err, "BuildPlugin: c.builderDo")
	}

	var t BuildPluginResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return BuildPluginResponse{}, errors.Wrap(err, "GetBuilderFeatures: dec.Decode")
	}

	return t, nil
}

type BuilderFeaturesResponse struct {
	Features  []string `json:"features"`
	Languages []struct {
		ID         string `json:"identifier"`
		ShortName  string `json:"short"`
		PrettyName string `json:"pretty"`
	}
}

func (c *Client2) GetBuilderFeatures(ctx context.Context) (BuilderFeaturesResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.host+pathBuilderFeaturesOld, nil)
	if err != nil {
		return BuilderFeaturesResponse{}, errors.Wrap(err, "GetBuilderFeatures: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return BuilderFeaturesResponse{}, errors.Wrap(err, "GetBuilderFeatures: c.do")
	}

	if res.StatusCode != http.StatusOK {
		return BuilderFeaturesResponse{}, fmt.Errorf("GetBuilderFeatures: received non-200 response %d", res.StatusCode)
	}

	// Marshal response body into what we need to give back.
	var t BuilderFeaturesResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return BuilderFeaturesResponse{}, errors.Wrap(err, "GetBuilderFeatures: dec.Decode")
	}

	return t, nil
}

type TestPluginDraftResponse struct{}

func (c *Client2) TestPluginDraft(ctx context.Context) (TestPluginDraftResponse, error) {
	return TestPluginDraftResponse{}, nil
}

// GetPluginDraft returns the currently set plugin draft for the given session token. To change the draft or the
// language you can use the CreatePluginDraft method instead with the name of a template.
func (c *Client2) GetPluginDraft(ctx context.Context, token CreateSessionResponse) (DraftResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.host+pathDraft, nil)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "GetPluginDraft: http.NewRequest")
	}

	res, err := c.builderDo(ctx, req, token)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "GetPluginDraft: c.builderDo")
	}

	b, _ := io.ReadAll(res.Body)

	fmt.Printf("all the body:\n%s\n", string(b))

	var t DraftResponse
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "GetPluginDraft: dec.Decode")
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
func (c *Client2) CreatePluginDraft(ctx context.Context, templateName string, token CreateSessionResponse) (DraftResponse, error) {
	if templateName == "" {
		return DraftResponse{}, errors.New("template name cannot be blank")
	}

	r := createDraftRequest{Template: templateName}
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(r)
	if err != nil {
		return DraftResponse{}, errors.Wrapf(err, "CreatePluginDraft: json.NewEncoder.Encode(createDraftRequest with template name '%s'", templateName)
	}

	req, err := http.NewRequest(http.MethodPost, c.host+pathDraft, &b)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "CreatePluginDraft: http.NewRequest")
	}

	res, err := c.builderDo(ctx, req, token)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "CreatePluginDraft: c.builderDo")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, _ := io.ReadAll(res.Body)
	fmt.Printf("returned body from create draft plugin\n\n%s\n\n", string(body))

	if res.StatusCode != http.StatusOK {
		return DraftResponse{}, fmt.Errorf("CreatePluginDraft: unexpected response code. Wanted %d, got %d", http.StatusOK, res.StatusCode)
	}

	var t DraftResponse
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return DraftResponse{}, errors.Wrap(err, "CreatePluginDraft: dec.Decode")
	}

	return t, nil
}

type PromotePluginDraftResponse struct{}

func (c *Client2) PromotePluginDraft(ctx context.Context) (PromotePluginDraftResponse, error) {
	return PromotePluginDraftResponse{}, nil
}

// builderDo is a common method to work with requests against the builder where a session token is needed instead of the
// environment token that the do method uses.
func (c *Client2) builderDo(ctx context.Context, req *http.Request, token CreateSessionResponse) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Add("Authorization", "Bearer "+token.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "c.builderDo: c.httpClient.Do")
	}

	return res, nil
}
