package se2

import "fmt"

type Module struct {
	Environment string
	Tenant      string
	Namespace   string
	Name        string
}

// NewModule instantiates a local v1.0.0 Module that can be used for various calls with compute.Client.
// Note: this constructor alone does not perform any actions on a remote Compute instance.
func NewModule(environment, tenant, namespace, name string) *Module {
	module := &Module{
		Environment: environment,
		Tenant:      tenant,
		Namespace:   namespace,
		Name:        name,
	}

	return module
}

func (m *Module) URI() string {
	return fmt.Sprintf("/%s.%s/%s/%s", m.Environment, m.Tenant, m.Namespace, m.Name)
}
