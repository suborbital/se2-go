package compute

import "github.com/pkg/errors"

type Client struct {
	config *Config

	adminAdapter     *adminAdapter
	builderAdapter   *builderAdapter
	executionAdapter *executionAdapter
}

func NewClient(config *Config) (*Client, error) {
	failErr := errors.New("failed to NewClient")

	if config == nil {
		return nil, errors.Wrap(failErr, "config cannot be nil")
	}

	client := &Client{
		config:           config,
		adminAdapter:     newAdminAdapter(config.adminConfig),
		builderAdapter:   newBuilderAdapter(config.builderConfig),
		executionAdapter: newExecutionAdapter(config.executionConfig),
	}

	return client, nil
}

func (c *Client) NewRunnable(environment, customerID, namespace, fnName string) *Runnable {
	runnable := &Runnable{
		client:       c,
		environment:  environment,
		customerID:   customerID,
		namespace:    namespace,
		functionName: fnName,
	}

	return runnable
}
