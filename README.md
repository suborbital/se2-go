# se2-go

Go client library for the Suborbital Extension Engine (SE2)

## Usage

In a Go project, run

```bash
go get github.com/suborbital/se2-go@latest
```

Every operation with SE2 is done with a `se2.Client`. Here's a simple example that fetches exisiting Plugins for a user and namespace.

```go
package main

import (
    "log"

    "github.com/suborbital/se2-go"
)

func main() {
    token, exists := os.LookupEnv("SE2_ENV_TOKEN")
    if !exists {
        log.Fatal("could not find token")
    }

    client, err := se2.NewClient(se2.LocalConfig(), token)
    if err != nil {
        log.Fatal(err)
    }

    // get a list of Plugins
    plugins, err := client.UserPlugins("userID", "namespace")
    if err != nil {
        log.Fatal(err)
    }

    for _, r := range plugins {
        log.Println(r.FQMN)
    }
}
```

See [examples](examples/) folder for more.
