# se2-go

Go client library for the Suborbital Extension Engine (SE2)

## Usage

In a Go project, run this to fetch the module:

```bash
go get github.com/suborbital/se2-go@latest
```

Every operation with SE2 is done with a `se2.Client`.

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/suborbital/se2-go"
)

func main() {
	token, exists := os.LookupEnv("SE2_ENV_TOKEN")
	if !exists {
		log.Fatal("could not find token")
	}

	client, err := se2.NewClient(se2.ModeProduction, token)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cxl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cxl()

	// get a list of plugins
	plugins, err := client.GetPlugins(ctx, "tenantName")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range plugins {
		log.Println(r.FQMN)
	}
}
```

Some of the methods require an API key that you get by creating a new access key in the admin area of https://suborbital.network. That API key is restricted to a specific environment. Every request is going to be against the environment that the key belongs to.

Others are session, or plugin specific. For those you need to first create a session, and then reuse the token with all the endpoints.

## Available methods
Each of these methods can be seen in the ["everything" annotated example](examples/everything).
### Tenants
#### GetTenantByName
```go
func example() {
    tenant, err := client.GetTenantByName(ctx, "tenantName")
	if err != nil {
		// handle error
    }
}
```
`GetTenantByName` will fetch the tenant by the given name. It requires a configured client with an API key that is valid. If the request was not successful, or the response was anything other than a 200 OK, it will return an error.

The returned tenant's structure:
```go
type TenantResponse struct {
	AuthorizedParty string `json:"authorized_party"`
	ID              string `json:"id"`
	Environment     string `json:"environment"`
	Name            string `json:"name"`
	Description     string `json:"description"`
}
```

#### CreteTenant

```go
func example() {
	tenant, err := client.CreateTenant(ctx, "tenantName", "tenant description")
    if err != nil {
        // handle error
    }
}
```
`CrateTenant` will create tenant with a given name and description. A successful response looks like `TenantResponse`. Requires a configured client and a correct API key. If the tenant wasn't created, this method returns a non nil error.

#### ListTenants

```go
func example() {
	tenants, err := client.ListTenants(ctx)
	if err != nil {
		// handle error
    }
}
```
`ListTenants` will list all the tenants within the environment that the access key belongs to. The response looks like this:
```go
type ListTenantResponse struct {
	Tenants []TenantResponse
}
```

#### UpdateTenantByName
`UpdateTenantByName` will update a tenant by the name you specified with a new description. No other parts of a tenant can be updated. It will return a non nil error if the update failed. This usually happens because the tenant was not found, or there was an authentication issue.

This method requires a configured client with a valid API key.
```go
func example() {
	tenant, err := client.UpdateTenantByName(ctx, "tenantName", "new description")
}
```

#### DeleteTenantByName

`DeleteTenantByName` will delete the tenant by the name if it finds it. Requires a configured client with a valid API key. If there was an issue deleting the tenant, it will return a non nil error. This is usually because a tenant by that name does not exist in the environment the API key specifies, or there was another authentication issue.

```go
func example() {
	err := client.DeleteTenantByName(ctx, "tenantName")
	if err != nil {
		// handle error
    }
}
```
