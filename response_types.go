package compute

import (
	"regexp"
	"time"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type RunnableResponse struct {
	Namespace    string `json:"namespace"`
	FunctionName string `json:"name"`
	Version      string `json:"version"`
	Language     string `json:"lang"`
	DraftVersion string `json:"draftVersion"`
	APIVersion   string `json:"apiVersion"`
	FQFN         string `json:"fqfn"`
}

func (r RunnableResponse) String() string {
	return r.FQFN
}

var fqfnRegexp = regexp.MustCompile(`(?P<environment>.*\..*)\.(?P<customer>.*)\#(?P<namespace>.*)::`)

// ToRunnable converts a RunnableResponse type to a Runnable
func (r RunnableResponse) ToRunnable() *Runnable {
	match := fqfnRegexp.FindStringSubmatch(r.FQFN)
	result := make(map[string]string)
	for i, name := range fqfnRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return &Runnable{
		environment:  result["environment"],
		customerID:   result["customer"],
		namespace:    result["namespace"],
		functionName: r.FunctionName,
		version:      r.Version,
		draftVersion: r.DraftVersion,
		language:     r.Language,
		apiVersion:   r.APIVersion,
	}
}

type UserFunctionsResponse struct {
	Functions []RunnableResponse `json:"functions"`
}

type ExecResult struct {
	UUID            string            `json:"uuid"`
	Timestamp       time.Time         `json:"timestamp"`
	Response        string            `json:"response"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
}

type ExecResultsResponse struct {
	Results []ExecResult `json:"results"`
}

type ExecErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"string"`
}
