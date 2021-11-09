package compute

import (
	"fmt"
)

type Runnable struct {
	environment  string
	customerID   string
	namespace    string
	functionName string
	version      string
	language     string
	draftVersion string
	apiVersion   string

	editorToken string
}

func (r Runnable) Environment() string {
	return r.environment
}

func (r Runnable) CustomerID() string {
	return r.customerID
}

func (r Runnable) Namespace() string {
	return r.namespace
}

func (r Runnable) FunctionName() string {
	return r.functionName
}

func (r Runnable) Version() string {
	return r.version
}

func (r Runnable) Path() string {
	return fmt.Sprintf("%s.%s/%s/%s", r.Environment(), r.CustomerID(), r.Namespace(), r.FunctionName())
}

func (r Runnable) String() string {
	return r.Path()
}

func (r Runnable) VersionPath() string {
	return fmt.Sprintf("%s/%s", r.Path(), r.Version())
}

func (r *Runnable) Token() string {
	return r.editorToken
}
