package compute

import (
	"time"

	"github.com/suborbital/atmo/directive"
)

type TokenResponse struct {
	Token string `json:"token"`
}
type UserFunctionsResponse struct {
	Functions []*directive.Runnable `json:"functions"`
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
