package compute_test

import (
	"net/http"
	"testing"

	"github.com/suborbital/compute-go"
)

func TestGetToken(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	runnable := client.NewRunnable(
		"com.suborbital",
		"customer",
		"default",
		"foo")

	_, err = client.EditorToken(runnable)
	if err != nil {
		t.Fatal(err)
	}

	expect := "L7rRBAgx8vcOtOJO2kBbjqrs"
	if token := runnable.Token(); token != expect {
		t.Fatalf("got %s, wanted %s", token, expect)
	}
}

func TestUserFunctions(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	fns, res, err := client.UserFunctions("customer", "default")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatal("expected 200, got", res.StatusCode)
	}

	t.Log(fns)
}

func TestGetAndExec(t *testing.T) {
	t.Parallel()

	client, err := compute.NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	fns, res, err := client.UserFunctions("customer", "default")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatal("expected 200, got", res.StatusCode)
	}

	if len(fns) < 1 {
		t.Skip("no runnables defined")
	}

	runnable := fns[0]

	result, _, err := client.ExecString(runnable, "world")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(result))

	// Tests the administrative results endpoint
	t.Run("ExecResults", func(t *testing.T) {
		execRes, _, err := client.FunctionExecResults(runnable)
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
}
