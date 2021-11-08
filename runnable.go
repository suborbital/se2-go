package compute

import (
	"fmt"
	"sync"
)

type Runnable struct {
	environment  string
	customerID   string
	namespace    string
	functionName string
	version      string

	token string

	lock sync.Mutex
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

func (r Runnable) VersionPath() string {
	return fmt.Sprintf("%s/%s", r.Path(), r.Version())
}

func (r *Runnable) setToken(token string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.token = token
}

func (r *Runnable) Token() string {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.token
}
