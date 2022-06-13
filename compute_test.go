package compute_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/suborbital/compute-go"
)

var environment = "com.suborbital"
var userID = ""
var envToken = ""

func TestMain(m *testing.M) {
	tok, exists := os.LookupEnv("SCC_ENV_TOKEN")
	if !exists {
		log.Fatal("SCC_ENV_TOKEN must be set to run tests!")
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

// TestBuilder must run before tests that depend on Runnables existing in SCN
func TestBuilder(t *testing.T) {
	client, err := compute.NewLocalClient(envToken)
	if err != nil {
		t.Fatal(err)
	}

	runnable := compute.NewRunnable(environment, userID, "default", "foo", "assemblyscript")

	t.Run("Template", func(t *testing.T) {
		template, err := client.BuilderTemplate(runnable)
		if err != nil {
			t.Fatal(err)
		}

		if template.Lang != runnable.Lang {
			t.Errorf("got Lang: '%s', want '%s'", template.Lang, runnable.Lang)
		}

		t.Logf("got template for '%s', length: %d", template.Lang, len(template.Contents))

		t.Run("Build", func(t *testing.T) {
			buildResult, err := client.BuildFunctionString(runnable, template.Contents)

			if err != nil {
				t.Fatal(err)
			}

			t.Log(buildResult)

			t.Run("GetDraft", func(t *testing.T) {
				editorState, err := client.GetDraft(runnable)
				if err != nil {
					t.Fatal(err)
				}

				if editorState.Contents != template.Contents {
					t.Error("function contents changed between build and draft")
				}
			})

			t.Run("Promote", func(t *testing.T) {
				promoteResponse, err := client.PromoteDraft(runnable)
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("runnable promoted: (%s -> %s)", runnable.Version, promoteResponse.Version)
			})
		})
	})
}

func TestUserFunctions(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient(envToken)
	if err != nil {
		t.Fatal(err)
	}

	identifier := fmt.Sprintf("%s.%s", environment, userID)

	fns, err := client.UserFunctions(identifier, "default")
	if err != nil {
		t.Fatal(err)
	}

	for _, fn := range fns {
		t.Log(fn.FQFN)
	}
}

func TestGetAndExec(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient(envToken)
	if err != nil {
		t.Fatal(err)
	}

	identifier := fmt.Sprintf("%s.%s", environment, userID)

	fns, err := client.UserFunctions(identifier, "default")
	if err != nil {
		t.Fatal(err)
	}

	if len(fns) < 1 {
		t.Skip("no runnables defined")
	}

	runnable := fns[0]

	result, uuid, err := client.ExecString(runnable, "world")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(result))

	// Give the results a moment to propagate
	time.Sleep(1 * time.Second)

	// Tests the administrative results endpoint
	t.Run("ExecResultsMetadata", func(t *testing.T) {
		execRes, err := client.FunctionResultsMetadata(runnable)
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
		res, err := client.FunctionResultMetadata(uuid)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("result:", res.UUID)
	})

	// Tests the administrative results endpoint
	t.Run("ExecResult", func(t *testing.T) {
		res, err := client.FunctionResult(uuid)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("result:", string(res))
	})
}

func TestBuilderHealth(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient(envToken)
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

	client, err := compute.NewLocalClient(envToken)
	if err != nil {
		t.Fatal(err)
	}

	features, err := client.BuilderFeatures()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(features)
}
