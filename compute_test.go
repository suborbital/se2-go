package compute_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/suborbital/compute-go"
	admin "github.com/suborbital/compute-go/compute/administrative"
	builder "github.com/suborbital/compute-go/compute/builder"
)

func TestAdministrativeDirect(t *testing.T) {
	t.Parallel()

	conf := admin.NewConfiguration()
	api := admin.NewAPIClient(conf)

	req := api.DefaultApi.GetToken(context.Background(),
		"com.suborbital",
		"customer",
		"default",
		"foo")

	tok, _, err := req.Execute()
	if err != nil {
		t.Fatal(err)
	}

	expect := "L7rRBAgx8vcOtOJO2kBbjqrs"
	if token := tok.GetToken(); token != expect {
		t.Fatalf("got %s, wanted %s", token, expect)
	}

	t.Log("token:", tok.GetToken())
}

func TestAdministrativePretty(t *testing.T) {
	t.Parallel()

	runnable, err := tokenBasic()
	if err != nil {
		t.FailNow()
	}

	expect := "L7rRBAgx8vcOtOJO2kBbjqrs"
	if token := runnable.Token(); token != expect {
		t.Fatalf("got %s, wanted %s", token, expect)
	}
}

func tokenBasic() (*compute.Runnable, error) {
	client, err := compute.NewClient(compute.DefaultConfig())
	if err != nil {
		return nil, err
	}

	runnable := client.NewRunnable(
		"com.suborbital",
		"customer",
		"default",
		"foo")

	return runnable, nil
}

func assemblyScriptFn() compute.Function {
	body, _ := ioutil.ReadFile("tests/hello.ts")

	return compute.Function{
		Language: "assemblyscript",
		Body:     string(body),
	}
}

func TestBuilderWithTokenDirect(t *testing.T) {
	t.Parallel()

	conf := builder.NewConfiguration()
	api := builder.NewAPIClient(conf)

	runnable, err := tokenBasic()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), builder.ContextAccessToken, runnable.Token())

	body, err := ioutil.ReadFile("tests/hello.ts")
	if err != nil {
		t.Fatal(err)
	}

	req := api.DefaultApi.BuildFunction(
		ctx,
		"assemblyscript",
		runnable.Environment(),
		runnable.CustomerID(),
		runnable.Namespace(), // weird method chaining, but okay
		runnable.FunctionName()).Body(string(body))

	output, res, err := req.Execute()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
	t.Log(res.Status)
}

func TestBuilderWithTokenPretty(t *testing.T) {
	t.Parallel()
	runnable, err := tokenBasic()
	if err != nil {
		t.Fatal(err)
	}

	output, res, err := runnable.BuildWith(assemblyScriptFn())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
	t.Log(res.Status)
}
