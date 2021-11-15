package compute_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/suborbital/compute-go"
)

var userID = ""

func TestMain(m *testing.M) {
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
	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	runnable := compute.NewRunnable("com.suborbital", userID, "default", "foo", "assemblyscript")

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

	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	fns, err := client.UserFunctions(userID, "default")
	if err != nil {
		t.Fatal(err)
	}

	for _, fn := range fns {
		t.Log(fn.FQFN)
	}
}

func TestGetAndExec(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	fns, err := client.UserFunctions(userID, "default")
	if err != nil {
		t.Fatal(err)
	}

	if len(fns) < 1 {
		t.Skip("no runnables defined")
	}

	runnable := fns[0]

	result, err := client.ExecString(runnable, "world")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(result))

	// Tests the administrative results endpoint
	t.Run("ExecResults", func(t *testing.T) {
		execRes, err := client.FunctionExecResults(runnable)
		if err != nil {
			t.Fatal(err)
		}

		if len(execRes.Results) < 1 {
			t.Fatal("expected at least one result")
		}
		sample := execRes.Results[0]
		t.Log("latest result:", sample.UUID, sample.Response)
		t.Logf("(%d total execution results)", len(execRes.Results))
	})

	// Tests the administrative results endpoint
	t.Run("ExecErrors", func(t *testing.T) {
		execErrs, err := client.FunctionExecErrors(runnable)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("errors:", execErrs.Errors)
		t.Logf("(%d total execution errors)", len(execErrs.Errors))
	})
}

func TestBuilderHealth(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient()
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

	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	features, err := client.BuilderFeatures()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(features.Features)
}
