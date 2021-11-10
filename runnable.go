package compute

import (
	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/atmo/fqfn"
)

// NewRunnableVersion instantiates a local versioned Runnable that can be used for various calls with compute.Client.
// Note: this constructor alone does not perform any actions on a remote Compute instance.
func NewRunnableVersion(environment, customerId, namespace, fnName, version, language string) *directive.Runnable {
	fqfnStr := fqfn.FromParts(environment+"."+customerId, namespace, fnName, version)
	FQFN := fqfn.Parse(fqfnStr)

	runnable := &directive.Runnable{
		Name:      fnName,
		Namespace: namespace,
		Lang:      language,
		Version:   version,
		FQFN:      fqfnStr,
		FQFNURI:   FQFN.HeadlessURLPath(),
	}

	return runnable
}

// NewRunnable instantiates a local v1.0.0 Runnable that can be used for various calls with compute.Client.
// Note: this constructor alone does not perform any actions on a remote Compute instance.
func NewRunnable(environment, customerId, namespace, fnName, language string) *directive.Runnable {
	return NewRunnableVersion(environment, customerId, namespace, fnName, "v1.0.0", language)
}
