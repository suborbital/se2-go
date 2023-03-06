package se2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const pathPlugins = pathTenantByName + "/plugins"

type Plugin2 struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Lang       string `json:"lang"`
	Ref        string `json:"ref"`
	ApiVersion string `json:"apiVersion"`
	Fqmn       string `json:"fqmn"`
	Uri        string `json:"uri"`
}

type PluginResponse struct {
	Plugins []Plugin2 `json:"plugins"`
}

func (c *Client2) GetPlugins(ctx context.Context, tenantName string) (PluginResponse, error) {
	if tenantName == "" {
		return PluginResponse{}, errors.New("tenant name cannot be blank")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(c.host+pathPlugins, tenantName), nil)
	if err != nil {
		return PluginResponse{}, errors.Wrap(err, "GetPlugins: http.NewRequest")
	}

	res, err := c.do(ctx, req)
	if err != nil {
		return PluginResponse{}, errors.Wrap(err, "GetPlugins: c.do")
	}

	var t PluginResponse

	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&t)
	if err != nil {
		return PluginResponse{}, errors.Wrap(err, "GetPlugins: dec.Decode")
	}

	return t, nil
}
