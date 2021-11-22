# compute-go
Go client library for Suborbital Compute

## Usage

In a Go project, run
```bash
go get github.com/suborbital/compute-go@latest
```

Every operation with Compute is done with a `compute.Client`. Here's a simple example that fetches exisiting Runnables for a user and namespace.

```go
package main

import (
    "log"

    "github.com/suborbital/compute-go"
)

func main() {
	token, err := os.LookupEnv("SCC_ENV_TOKEN")
    if err != nil {
        log.Fatal(err)
    }

	client, err := compute.NewClient(compute.LocalConfig(), token)
    if err != nil {
        log.Fatal(err)
    }

    // get a list of Runnables
    runnables, err := client.UserFunctions("userID", "namespace")
    if err != nil {
        log.Fatal(err)
    }

    for _, r := range runnables {
        log.Println(r.FQFN)
    }
}
```

See [examples](examples/) folder for more.
