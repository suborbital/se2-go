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
