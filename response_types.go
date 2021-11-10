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

type ExecError struct {
	UUID      string    `json:"uuid,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Code      int       `json:"code"`
	Message   string    `json:"string"`
}

type ExecErrorResponse struct {
	Errors []ExecError `json:"errors"`
}

type FeaturesResponse struct {
	Features []string `json:"features"`
}

// TestPayload is a single test for a Runnable
type TestPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Payload     string `json:"payload"`
}

// EditorStateResponse is a response to requests to get editorState
type EditorStateResponse struct {
	Lang     string        `json:"lang"`
	Contents string        `json:"contents"`
	Tests    []TestPayload `json:"tests"`
}

type BuildResult struct {
	Succeeded bool   `json:"succeeded"`
	OutputLog string `json:"outputLog"`
}
