package compute

import (
	"context"
	"net/http"

	admin "github.com/suborbital/compute-go/compute/administrative"
)

type adminAdapter struct {
	client *admin.APIClient
	ctx    context.Context
}

func newAdminAdapter(conf *admin.Configuration) *adminAdapter {
	adapter := &adminAdapter{
		client: admin.NewAPIClient(conf),
		ctx:    context.Background(),
	}

	return adapter
}

func (a adminAdapter) GetToken(runnable *Runnable) (string, *http.Response, error) {
	req := a.client.DefaultApi.GetToken(a.ctx,
		runnable.Environment(),
		runnable.CustomerID(),
		runnable.Namespace(),
		runnable.FunctionName())

	token, res, err := req.Execute()
	return token.GetToken(), res, err
}
