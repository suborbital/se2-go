package se2_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/suborbital/se2-go"
)

var environment = "com.suborbital"
var userID = ""
var envToken = ""

var tmpl = "tinygo"

func TestMain(m *testing.M) {
	tok, exists := os.LookupEnv("SE2_ENV_TOKEN")
	if !exists {
		log.Fatal("SE2_ENV_TOKEN must be set to run tests!")
	}
	envToken = tok

	// creates a new user for every test run
	uuid := uuid.New()
	raw, _ := uuid.MarshalBinary()

	userID = base64.URLEncoding.EncodeToString(raw)

	// chop off base64 ==
	userID = userID[:len(userID)-2]

	os.Exit(m.Run())
}

func TestUserID(t *testing.T) {
	t.Logf("Using UserID: %s", userID)
}

// TestBuilder must run before tests that depend on plugins existing in SE2
func TestBuilder(t *testing.T) {
	client, err := se2.NewLocalClient(envToken)
	if err != nil {
		t.Fatal(err)
	}

	plugin := se2.NewPlugin(environment, userID, "default", "foo")

	t.Run("Template", func(t *testing.T) {
		template, err := client.BuilderTemplate(plugin, tmpl)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("got template for '%s', length: %d", template.Lang, len(template.Contents))

		t.Run("Build", func(t *testing.T) {
			buildResult, err := client.BuildPluginString(plugin, tmpl, template.Contents)

			if err != nil {
				t.Fatal(err)
			}

			t.Log(buildResult)

			t.Run("GetDraft", func(t *testing.T) {
				editorState, err := client.GetDraft(plugin)
				if err != nil {
					t.Fatal(err)
				}

				if editorState.Contents != template.Contents {
					t.Error("plugin contents changed between build and draft")
				}
			})

			t.Run("Promote", func(t *testing.T) {
				promoteResponse, err := client.PromoteDraft(plugin)
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("plugin promoted to %s", promoteResponse.Version)
				time.Sleep(time.Second * 2)
			})
		})
	})
}

func TestUserPlugins(t *testing.T) {
	t.Parallel()

	client, err := se2.NewLocalClient(envToken)
	if err != nil {
		t.Fatal(err)
	}

	identifier := fmt.Sprintf("%s.%s", environment, userID)

	fns, err := client.UserPlugins(identifier, "default")
	if err != nil {
		t.Fatal(err)
	}

	for _, fn := range fns {
		t.Log(fn.FQMN)
	}
}

func TestGetAndExec(t *testing.T) {
	t.Parallel()

	client, err := se2.NewLocalClient(envToken)
	if err != nil {
		t.Fatal(err)
	}

	identifier := fmt.Sprintf("%s.%s", environment, userID)

	fns, err := client.UserPlugins(identifier, "default")
	if err != nil {
		t.Fatal(err)
	}

	if len(fns) < 1 {
		t.Skip("no plugins defined")
	}

	plugin := se2.Plugin{
		Environment: environment,
		Tenant:      userID,
		Namespace:   fns[0].Namespace,
		Name:        fns[0].Name,
	}

	result, uuid, err := client.ExecString(&plugin, "world")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(result))

	// Give the results a moment to propagate
	time.Sleep(1 * time.Second)

	// Tests the administrative results endpoint
	t.Run("ExecResultsMetadata", func(t *testing.T) {
		execRes, err := client.ResultsMetadata(&plugin)
		if err != nil {
			t.Fatal(err)
		}

		if len(execRes) < 1 {
			t.Fatal("expected at least one result")
		}
		sample := execRes[0]
		t.Log("latest result:", sample.UUID)
		t.Logf("(%d total execution results)", len(execRes))
	})

	// Tests the administrative results endpoint
	t.Run("ExecResultMetadata", func(t *testing.T) {
		res, err := client.ResultMetadata(uuid)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("result:", res.UUID)
	})

	// Tests the administrative results endpoint
	t.Run("ExecResult", func(t *testing.T) {
		res, err := client.ExecutionResult(uuid)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("result:", string(res))
	})
}

func TestBuilderHealth(t *testing.T) {
	t.Parallel()

	client, err := se2.NewLocalClient(envToken)
	if err != nil {
		t.Fatal(err)
	}

	healthy, err := client.BuilderHealth()
	if err != nil {
		t.Fatal(err)
	}

	if !healthy {
		t.Fatal("BuilderHealth returned false, want true")
	}
}

func TestBuilderFeatures(t *testing.T) {
	t.Parallel()

	client, err := se2.NewLocalClient(envToken)
	if err != nil {
		t.Fatal(err)
	}

	features, err := client.BuilderFeatures()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(features)
}
