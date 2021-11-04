# \DefaultApi

All URIs are relative to *http://local.suborbital.network:8082*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BuildFunction**](DefaultApi.md#BuildFunction) | **Post** /api/v1/build/{language}/{environment}.{customerID}/{namespace}/{fnName} | Builds the provided code using the specified language toolchain
[**DeployDraft**](DefaultApi.md#DeployDraft) | **Post** /api/v1/draft/{environment}.{customerId}/{namespace}/{fnName}/promote | Deploys the specified runnable
[**GetDraft**](DefaultApi.md#GetDraft) | **Get** /api/v1/draft/{environment}.{customerId}/{namespace}/{fnName} | Gets the draft for the specified runnable
[**GetFeatures**](DefaultApi.md#GetFeatures) | **Get** /api/v1/features | Returns a list of supported builder features. Primarily used to detect the presence of a testing service.
[**GetHealth**](DefaultApi.md#GetHealth) | **Get** /api/v1/health | Returns an OK response to indicate a healthy service (returns no body)
[**GetTemplate**](DefaultApi.md#GetTemplate) | **Get** /api/v1/template/{language} | Gets the template for a new function of the given language
[**GetTemplateV2**](DefaultApi.md#GetTemplateV2) | **Get** /api/v2/template/{language}/{fnName} | Gets the template for a new function of the given language initialized with the supplied function name
[**TestDraft**](DefaultApi.md#TestDraft) | **Post** /api/v1/test/{environment}.{customerId}/{namespace}/{fnName} | Tests drafts



## BuildFunction

> string BuildFunction(ctx, language, environment, customerID, namespace, fnName).Body(body).Execute()

Builds the provided code using the specified language toolchain

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    language := "assemblyscript" // string | The language toolchain used to configure the builder
    environment := "com.suborbital" // string | The root compute environment (i.e. the vendor)
    customerID := "acmeco" // string | The vendor's customer (i.e. the user)
    namespace := "default" // string | The function namespace (vendor-defined groups of functions)
    fnName := "httpget" // string | The function name (customer-defined)
    body := "body_example" // string | Bytes of the code to be built, for example the contents of lib.rs or lib.ts

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.BuildFunction(context.Background(), language, environment, customerID, namespace, fnName).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.BuildFunction``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `BuildFunction`: string
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.BuildFunction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**language** | **string** | The language toolchain used to configure the builder | 
**environment** | **string** | The root compute environment (i.e. the vendor) | 
**customerID** | **string** | The vendor&#39;s customer (i.e. the user) | 
**namespace** | **string** | The function namespace (vendor-defined groups of functions) | 
**fnName** | **string** | The function name (customer-defined) | 

### Other Parameters

Other parameters are passed through a pointer to a apiBuildFunctionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





 **body** | **string** | Bytes of the code to be built, for example the contents of lib.rs or lib.ts | 

### Return type

**string**

### Authorization

[bearerAuthentication](../README.md#bearerAuthentication)

### HTTP request headers

- **Content-Type**: text/plain
- **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeployDraft

> string DeployDraft(ctx, environment, customerID, namespace, fnName).Execute()

Deploys the specified runnable

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    environment := "com.suborbital" // string | The root compute environment (i.e. the vendor)
    customerID := "acmeco" // string | The vendor's customer (i.e. the user)
    namespace := "default" // string | The function namespace (vendor-defined groups of functions)
    fnName := "httpget" // string | The function name (customer-defined)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.DeployDraft(context.Background(), environment, customerID, namespace, fnName).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeployDraft``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeployDraft`: string
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.DeployDraft`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**environment** | **string** | The root compute environment (i.e. the vendor) | 
**customerID** | **string** | The vendor&#39;s customer (i.e. the user) | 
**namespace** | **string** | The function namespace (vendor-defined groups of functions) | 
**fnName** | **string** | The function name (customer-defined) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeployDraftRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





### Return type

**string**

### Authorization

[bearerAuthentication](../README.md#bearerAuthentication)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDraft

> EditorState GetDraft(ctx, environment, customerID, namespace, fnName).Execute()

Gets the draft for the specified runnable

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    environment := "com.suborbital" // string | The root compute environment (i.e. the vendor)
    customerID := "acmeco" // string | The vendor's customer (i.e. the user)
    namespace := "default" // string | The function namespace (vendor-defined groups of functions)
    fnName := "httpget" // string | The function name (customer-defined)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.GetDraft(context.Background(), environment, customerID, namespace, fnName).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetDraft``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetDraft`: EditorState
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetDraft`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**environment** | **string** | The root compute environment (i.e. the vendor) | 
**customerID** | **string** | The vendor&#39;s customer (i.e. the user) | 
**namespace** | **string** | The function namespace (vendor-defined groups of functions) | 
**fnName** | **string** | The function name (customer-defined) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDraftRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





### Return type

[**EditorState**](EditorState.md)

### Authorization

[bearerAuthentication](../README.md#bearerAuthentication)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetFeatures

> InlineResponse200 GetFeatures(ctx).Execute()

Returns a list of supported builder features. Primarily used to detect the presence of a testing service.

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.GetFeatures(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetFeatures``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetFeatures`: InlineResponse200
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetFeatures`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetFeaturesRequest struct via the builder pattern


### Return type

[**InlineResponse200**](InlineResponse200.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetHealth

> GetHealth(ctx).Execute()

Returns an OK response to indicate a healthy service (returns no body)

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.GetHealth(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetHealth``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetHealthRequest struct via the builder pattern


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTemplate

> EditorState GetTemplate(ctx, language).Namespace(namespace).Execute()

Gets the template for a new function of the given language

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    language := "assemblyscript" // string | The language toolchain used to configure the builder
    namespace := "namespace_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.GetTemplate(context.Background(), language).Namespace(namespace).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetTemplate``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetTemplate`: EditorState
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetTemplate`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**language** | **string** | The language toolchain used to configure the builder | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **namespace** | **string** |  | 

### Return type

[**EditorState**](EditorState.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTemplateV2

> EditorState GetTemplateV2(ctx, language, fnName).Namespace(namespace).Execute()

Gets the template for a new function of the given language initialized with the supplied function name

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    language := "assemblyscript" // string | The language toolchain used to configure the builder
    fnName := "httpget" // string | The function name (customer-defined)
    namespace := "namespace_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.GetTemplateV2(context.Background(), language, fnName).Namespace(namespace).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetTemplateV2``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetTemplateV2`: EditorState
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetTemplateV2`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**language** | **string** | The language toolchain used to configure the builder | 
**fnName** | **string** | The function name (customer-defined) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTemplateV2Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **namespace** | **string** |  | 

### Return type

[**EditorState**](EditorState.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TestDraft

> InlineResponse2001 TestDraft(ctx, environment, customerID, namespace, fnName).Execute()

Tests drafts

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    environment := "com.suborbital" // string | The root compute environment (i.e. the vendor)
    customerID := "acmeco" // string | The vendor's customer (i.e. the user)
    namespace := "default" // string | The function namespace (vendor-defined groups of functions)
    fnName := "httpget" // string | The function name (customer-defined)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.TestDraft(context.Background(), environment, customerID, namespace, fnName).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.TestDraft``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `TestDraft`: InlineResponse2001
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.TestDraft`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**environment** | **string** | The root compute environment (i.e. the vendor) | 
**customerID** | **string** | The vendor&#39;s customer (i.e. the user) | 
**namespace** | **string** | The function namespace (vendor-defined groups of functions) | 
**fnName** | **string** | The function name (customer-defined) | 

### Other Parameters

Other parameters are passed through a pointer to a apiTestDraftRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





### Return type

[**InlineResponse2001**](InlineResponse2001.md)

### Authorization

[bearerAuthentication](../README.md#bearerAuthentication)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

