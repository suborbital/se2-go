package compute_test

import (
	"testing"

	"github.com/suborbital/compute-go"
)

func TestGetToken(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	runnable := compute.NewRunnable("com.suborbital", "customer", "default", "foo", "assemblyscript")

	token, err := client.EditorToken(runnable)
	if err != nil {
		t.Fatal(err)
	}

	expect := "L7rRBAgx8vcOtOJO2kBbjqrs"
	if token != expect {
		t.Fatalf("got %s, wanted %s", token, expect)
	}
}

func TestUserFunctions(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	fns, err := client.UserFunctions("customer", "default")
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

	fns, err := client.UserFunctions("customer", "default")
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

func TestBuilder(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	lang := "assemblyscript"
	namespace := "default"

	t.Run("Template/V1", func(t *testing.T) {
		template, err := client.BuilderTemplateV1(lang, namespace)
		if err != nil {
			t.Fatal(err)
		}

		if template.Lang != lang {
			t.Errorf("got Lang: '%s', want '%s'", template.Lang, lang)
		}

		t.Logf("got template for '%s', length: %d", template.Lang, len(template.Contents))
	})

	t.Run("Template/V2", func(t *testing.T) {
		template, err := client.BuilderTemplateV2(lang, namespace, "foo")
		if err != nil {
			t.Fatal(err)
		}

		if template.Lang != lang {
			t.Errorf("got Lang: '%s', want '%s'", template.Lang, lang)
		}

		t.Logf("got template for '%s', length: %d", template.Lang, len(template.Contents))

		t.Run("Build", func(t *testing.T) {
			runnable := compute.NewRunnable("com.suborbital", "customer", "default", "foo", "assemblyscript")
			buildResult, err := client.BuildFunctionString(runnable, template.Contents)

			if err != nil {
				t.Fatal(err)
			}

			t.Log(buildResult)
		})
	})
}
