package se2

import "fmt"

type Plugin struct {
	Environment string
	Tenant      string
	Namespace   string
	Name        string
}

// NewPlugin instantiates a local v1.0.0 plugin that can be used for various calls with se2.Client.
// Note: this constructor alone does not perform any actions on a remote SE2 instance.
func NewPlugin(environment, tenant, namespace, name string) *Plugin {
	plugin := &Plugin{
		Environment: environment,
		Tenant:      tenant,
		Namespace:   namespace,
		Name:        name,
	}

	return plugin
}

func (m *Plugin) URI() string {
	return fmt.Sprintf("/%s.%s/%s/%s", m.Environment, m.Tenant, m.Namespace, m.Name)
}