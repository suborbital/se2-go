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
)

type BuildPluginRequest struct{}

type BuildPluginResponse struct{}

func (c *Client2) BuildPlugin(ctx context.Context) (BuildPluginResponse, error) {
	return BuildPluginResponse{}, nil
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
	req, err := http.NewRequest(http.MethodGet, c.builderHost+pathBuilderFeaturesOld, nil)
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

type GetPluginDraftResponse struct{}

func (c *Client2) GetPluginDraft(ctx context.Context, token CreateSessionResponse) (GetPluginDraftResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.builderHost+pathDraft, nil)
	if err != nil {
		return GetPluginDraftResponse{}, errors.Wrap(err, "GetPluginDraft: http.NewRequest")
	}

	res, err := c.builderDo(ctx, req, token)
	if err != nil {
		return GetPluginDraftResponse{}, errors.Wrap(err, "GetPluginDraft: c.builderDo")
	}

	b, _ := io.ReadAll(res.Body)

	fmt.Printf("all the body:\n%s\n", string(b))

	var t GetPluginDraftResponse
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return GetPluginDraftResponse{}, errors.Wrap(err, "GetPluginDraft: dec.Decode")
	}

	return t, nil
}

type PromotePluginDraftResponse struct{}

func (c *Client2) PromotePluginDraft(ctx context.Context) (PromotePluginDraftResponse, error) {
	return PromotePluginDraftResponse{}, nil
}

func (c *Client2) builderDo(ctx context.Context, req *http.Request, token CreateSessionResponse) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Add("Authorization", "Bearer "+token.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "c.builderDo: c.httpClient.Do")
	}

	return res, nil
}