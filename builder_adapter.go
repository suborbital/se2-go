package compute

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	builder "github.com/suborbital/compute-go/compute/builder"
)

type builderAdapter struct {
	client *builder.APIClient
	ctx    context.Context
}

func newBuilderAdapter(conf *builder.Configuration) *builderAdapter {
	adapter := &builderAdapter{
		client: builder.NewAPIClient(conf),
		ctx:    context.Background(),
	}

	return adapter
}

func (a builderAdapter) BuildFunction(runnable *Runnable, fn Function) (string, *http.Response, error) {
	if runnable.Token() == "" {
		return "", nil, errors.New("could not get token for Runnable")
	}

	ctx := context.WithValue(a.ctx, builder.ContextAccessToken, runnable.Token())

	req := a.client.DefaultApi.BuildFunction(
		ctx,
		fn.Language,
		runnable.Environment(),
		runnable.CustomerID(),
		runnable.Namespace(),
		runnable.FunctionName()).Body(fn.Body)

	return req.Execute()
}
