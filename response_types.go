package se2

import (
	"time"

	"github.com/suborbital/systemspec/tenant"
)

type TokenResponse struct {
	Token string `json:"token"`
}
type UserFunctionsResponse struct {
	Functions []*tenant.Module `json:"functions"`
}

type ExecError struct {
	Code    int    `json:"code"`
	Message string `json:"string"`
}

type ExecMetadata struct {
	UUID      string    `json:"uuid"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
	Error     ExecError `json:"error"`
}

type FeatureLanguage struct {
	ID     string `json:"identifier"`
	Short  string `json:"short"`
	Pretty string `json:"pretty"`
}

type FeaturesResponse struct {
	Features  []string          `json:"features"`
	Langauges []FeatureLanguage `json:"languages"`
}

// TestPayload is a single test for a Module
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

type PromoteDraftResponse struct {
	Version string `json:"version"`
}
