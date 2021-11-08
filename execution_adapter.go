package compute

import (
	"context"

	execution "github.com/suborbital/compute-go/compute/execution"
)

type executionAdapter struct {
	client *execution.APIClient
	ctx    context.Context
}

func newExecutionAdapter(conf *execution.Configuration) *executionAdapter {
	adapter := &executionAdapter{
		client: execution.NewAPIClient(conf),
		ctx:    context.Background(),
	}

	return adapter
}
