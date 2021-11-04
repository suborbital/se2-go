package compute

import "net/http"

type Runnable struct {
	client       *Client
	environment  string
	customerID   string
	namespace    string
	functionName string

	token string
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

func (r *Runnable) Token() string {
	if r.token == "" {
		tok, _, err := r.client.adminAdapter.GetToken(r)
		if err == nil {
			r.token = tok
		}
	}
	return r.token
}

func (r *Runnable) BuildWith(fn Function) (string, *http.Response, error) {
	return r.client.builderAdapter.BuildFunction(r, fn)
}
