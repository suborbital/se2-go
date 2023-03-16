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

### Tenant methods

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

### Template methods

#### ListTemplates
`ListTemplates` will list all the available templates for the environment specified by the API key. Requires a configured client and a valid API key.
```go
func example() {
	templates, err := client.ListTemplates(ctx)
	if err != nil {
		// handle error
    }
}
```
Returned payload looks like this:
```go
type ListTemplatesResponse struct {
	Templates []Template `json:"templates"`
}

type Template struct {
    Name    string `json:"name"`
    Lang    string `json:"lang"`
    Main    string `json:"main,omitempty"`
    Version string `json:"api_version"`
}
```

#### GetTemplate

`GetTemplate` will return a single template specified by its name in the environment the API key is for. It requires a configured client and a valid API key.

```go
func example() {
	template, err := client.GetTemplate(ctx, "templateName")
	if err != nil {
		// handle error
    }
}
```

#### ImportTemplatesFromGitHub
`ImportTemplatesFromGitHub` will import templates found in a GitHub repository. It takes a repository name in the form of `organization/repository`, a reference, which can be a branch, tag, or commit sha, and a path, which is where the templates are within the repository.

To help with figuring out what either of those should be, the arguments will be substituted into this download link:
```go
https://github.com/{repository}/archive/{ref}.tar.gz
```
Any two values that would result in a 200 OK and a download response from GitHub will be valid arguments to the method.

This requires a configured client, a valid API key, and that the repository is public and not archived.

```go
func example() {
	err := client.ImportTemplatesFromGitHub(ctx, "suborbital/sdk", "main", "templates")
	if err != nil {
		// handle error
    }
}
```

### Plugin methods

#### GetPlugins

`GetPlugins` will list the plugins for a given tenant, which is scoped to the environment the API key specifies. Requires a configured client and a valid API key.

```go
func example() {
	plugins, err := client.GetPlugins(ctx, "tenantName")
    if err != nil {
		// handle error
    }
}
```
The response shape is a `PluginResponse`, which looks like this:
```go
type Plugin struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Lang       string `json:"lang"`
	Ref        string `json:"ref"`
	APIVersion string `json:"apiVersion"`
	FQMN       string `json:"fqmn"`
	URI        string `json:"uri"`
}

type PluginResponse struct {
	Plugins []Plugin `json:"plugins"`
}
```

### Session methods

There's only a single endpoint here. Use this to create a session to use with the builder methods.

#### CreateSession

`CreateSession` creates a session for a given tenant, namespace, and plugin name. Once you have the token, every request to the builder with the same session will relate to the same tenant-namespace-plugin.

This requires a configured client with a valid API key.

```go
func example() {
	session, err := client.CreateSession(ctx, "tenantName", "namespace", "pluginName")
    if err != nil {
        // handle error
    }
}
```

### Builder methods

#### GetBuilderFeatures
`GetBuilderFeatures` returns a list of things that the builder can do for you. This endpoint requires only a configured client and a valid API key. It does NOT need a session token.

```go
func example() {
	features, err := client.GetBuilderFeatures(ctx)
	
}
```
The shape of `features` is this:

```go
type BuilderFeaturesResponse struct {
    Features  []string `json:"features"`
    Languages []Languages
}

type Languages struct {
    ID         string `json:"identifier"`
    ShortName  string `json:"short"`
    PrettyName string `json:"pretty"`
}
```

#### CreatePluginDraft

`CreatePluginDraft` creates a plugin draft based on a template. You can query the list of available templates with the [`ListTemplates`](#listtemplates) call, or add new ones using the [`ImportTemplatesFromGitHub`](#importtemplatesfromgithub) method.

This endpoint requires a [session token](#createsession) to work. Pass in the variable as is you get from the `CreateSession` method.

```go
func example() {
	session, err := client.CreateSession(ctx, "tenantName", "namespace", "pluginName")
	if err != nil {
		// handle error
	}
	
	draft, err := client.CreatePluginDraft(ctx, "javascript", session)
	if err != nil {
		// handle error
    }
}
```
The returned `draft` variable will be of this structure, where `Lang` is the programming language, and `Contents` is the actual code of the starter state of the plugin.
```go
type DraftResponse struct {
	Lang     string `json:"lang"`
	Contents string `json:"contents"`
}
```

#### GetPluginDraft

`GetPluginDraft` returns the current state of the draft plugin specified by the session token. The response is the same as you get from [`CreatePluginDraft`](#createplugindraft). It requires a [session token](#createsession).

```go
func example() {
	draft, err := client.GetPluginDraft(ctx, session)
    if err != nil {
        // handle error
    }
}
```

#### BuildPlugin

`BuildPlugin` takes in a new plugin code body as a byte slice and a session token, and builds the plugin. It returns a non nil error if something went wrong. This is most likely going to be a syntax error in the supplied code, or a form of authentication error. It requires a non empty plugin code body, and a [session token](#createsession).

```go
func example() {
	output, err := client.BuildPlugin(ctx, pluginBody, session)
    if err != nil {
        // handle error
    }
}
```
The returned `output` has the following structure:

```go
type BuildPluginResponse struct {
	Succeeded bool   `json:"succeeded"`
	OutputLog string `json:"outputLog"`
}
```
The `OutputLog` is what the compiler printed to the terminal on the server. You can use it to debug what happened if the build did not succeed.

If `succeeded` is true, a new call to `GetPluginDraft` with the same session token will be the same code that you passed into the `BuildPlugin` method.

#### TestPluginDraft

`TestPluginDraft` takes in a byte slice to use as input to the current state of the plugin in the session, and will return a test response which has a string result for output, and an error if something went wrong while executing the plugin with the input.

```go
func example() {
	result, err := client.TestPluginDraft(ctx, []byte(`hello`), session)
    if err != nil {
        // handle error
    }
}
```
The structure of the returned `result` is this:
```go
type TestPluginDraftResponse struct {
	Result string   `json:"result"`
	Error  runError `json:"error"`
}

type runError struct {
    Code    int    `json:"code,omitempty"`
    Message string `json:"message,omitempty"`
}
```

#### PromotePluginDraft

`PromotePluginDraft` will push the current draft version to the live edge servers. It will take a short time while code propagates, but after this executing the plugin on the edge will use the new code. This requires a [session token](#createsession).

```go
func example() {
	response, err := client.PromotePluginDraft(ctx, session)
    if err != nil {
        // handle error
    }
}
```
If all went well, the structure of the returned `response` looks like this:
```go
type PromotePluginDraftResponse struct {
    Ref string `json:"ref"`
}
```


### Execution
Contains a single method to execute published, or promoted, plugins.

#### Exec
`Exec` takes a byte slice payload that will be used as input for the plugin, and a trio of ident, namespace, and plugin name to specify which plugin to execute. The trio of inputs is the same as you used in the [`CreateSession`](#createsession) call. If execution failed, `err` is going to be non nil.

This endpoint requires a configured client, and a valid API key, not a session token!

```go
func example() {
	responseBytes, err := client.Exec(ctx, []byte(`hello`), "tenantName", "namespace", "pluginName")
    if err != nil {
        // handle error
    }
}
```
`responseBytes` is a `[]byte` type. This is a bytes in, bytes out operation.
