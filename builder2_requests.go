package se2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	pathBuilderPrefix   = "/builder/b1"
	pathBuilderFeatures = pathBuilderPrefix + "/features"
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
	req, err := http.NewRequest(http.MethodGet, c.host+pathBuilderFeatures, nil)
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

func (c *Client2) GetPluginDraft(ctx context.Context) (GetPluginDraftResponse, error) {
	return GetPluginDraftResponse{}, nil
}

type PromotePluginDraftResponse struct{}

func (c *Client2) PromotePluginDraft(ctx context.Context) (PromotePluginDraftResponse, error) {
	return PromotePluginDraftResponse{}, nil
}
