# compute-go
Go client library for Suborbital Compute

## Usage

In a Go project, run
```bash
go get github.com/suborbital/compute-go
```

Every operation with Compute is done with a `compute.Client`. Here's a simple example that fetches exisiting Runnables for a customer and namespace.

```go
package main

import (
    "log"

    "github.com/suborbital/compute-go"
)

func main() {
    client, err := compute.NewLocalClient()
    if err != nil {
        log.Fatal(err)
    }

    // get a list of Runnables
    runnables, err := client.UserFunctions("customerID", "namespace")
    if err != nil {
        log.Fatal(err)
    }

    for _, r := range runnables {
        log.Println(r.FQFN)
    }
}
```

See [examples](examples/) folder for more.
