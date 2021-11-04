# Go API client for openapi

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 1.0
- Package version: 1.0.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/oauth2
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```golang
import sw "./openapi"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```golang
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `sw.ContextServerIndex` of type `int`.

```golang
ctx := context.WithValue(context.Background(), sw.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `sw.ContextServerVariables` of type `map[string]string`.

```golang
ctx := context.WithValue(context.Background(), sw.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```
ctx := context.WithValue(context.Background(), sw.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), sw.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *http://local.suborbital.network:8082*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*DefaultApi* | [**BuildFunction**](docs/DefaultApi.md#buildfunction) | **Post** /api/v1/build/{language}/{environment}.{customerID}/{namespace}/{fnName} | Builds the provided code using the specified language toolchain
*DefaultApi* | [**DeployDraft**](docs/DefaultApi.md#deploydraft) | **Post** /api/v1/draft/{environment}.{customerId}/{namespace}/{fnName}/promote | Deploys the specified runnable
*DefaultApi* | [**GetDraft**](docs/DefaultApi.md#getdraft) | **Get** /api/v1/draft/{environment}.{customerId}/{namespace}/{fnName} | Gets the draft for the specified runnable
*DefaultApi* | [**GetFeatures**](docs/DefaultApi.md#getfeatures) | **Get** /api/v1/features | Returns a list of supported builder features. Primarily used to detect the presence of a testing service.
*DefaultApi* | [**GetHealth**](docs/DefaultApi.md#gethealth) | **Get** /api/v1/health | Returns an OK response to indicate a healthy service (returns no body)
*DefaultApi* | [**GetTemplate**](docs/DefaultApi.md#gettemplate) | **Get** /api/v1/template/{language} | Gets the template for a new function of the given language
*DefaultApi* | [**GetTemplateV2**](docs/DefaultApi.md#gettemplatev2) | **Get** /api/v2/template/{language}/{fnName} | Gets the template for a new function of the given language initialized with the supplied function name
*DefaultApi* | [**TestDraft**](docs/DefaultApi.md#testdraft) | **Post** /api/v1/test/{environment}.{customerId}/{namespace}/{fnName} | Tests drafts


## Documentation For Models

 - [EditorState](docs/EditorState.md)
 - [InlineResponse200](docs/InlineResponse200.md)
 - [InlineResponse2001](docs/InlineResponse2001.md)
 - [TestPayload](docs/TestPayload.md)


## Documentation For Authorization



### bearerAuthentication

- **Type**: HTTP Bearer token authentication

Example

```golang
auth := context.WithValue(context.Background(), sw.ContextAccessToken, "BEARERTOKENSTRING")
r, err := client.Service.Operation(auth, args)
```


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author


