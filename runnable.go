package compute

import (
	"github.com/suborbital/systemspec/fqmn"
	tenantPkg "github.com/suborbital/systemspec/tenant"
)

func NewDraft() {

}

// NewModule instantiates a local v1.0.0 Runnable that can be used for various calls with compute.Client.
// Note: this constructor alone does not perform any actions on a remote Compute instance.
func NewModule(tenant, namespace, module, ref, language string) (*tenantPkg.Module, error) {
	fqmnStr, err := fqmn.FromParts(tenant, namespace, module, ref)
	if err != nil {
		return nil, err
	}

	runnable := &tenantPkg.Module{
		Name:      module,
		Namespace: namespace,
		Lang:      language,
		FQMN:      fqmnStr,
	}

	return runnable, nil
}
