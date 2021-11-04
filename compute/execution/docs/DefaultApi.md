# \DefaultApi

All URIs are relative to *http://local.suborbital.network:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RunFunction**](DefaultApi.md#RunFunction) | **Post** /{environment}.{customerID}/{namespace}/{fnName}/{versionNumber} | Executes the given function, with the provided body, params and state loaded into the function at runtime.



## RunFunction

> RunFunction(ctx, environment, customerID, namespace, fnName, versionNumber).Body(body).XAtmoState(xAtmoState).Execute()

Executes the given function, with the provided body, params and state loaded into the function at runtime.

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
    versionNumber := "v1.0.0" // string | The module version
    body := "body_example" // string | The payload bytes to be used as input to the function.
    xAtmoState := "xAtmoState_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.RunFunction(context.Background(), environment, customerID, namespace, fnName, versionNumber).Body(body).XAtmoState(xAtmoState).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.RunFunction``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
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
**versionNumber** | **string** | The module version | 

### Other Parameters

Other parameters are passed through a pointer to a apiRunFunctionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





 **body** | **string** | The payload bytes to be used as input to the function. | 
 **xAtmoState** | **string** |  | 

### Return type

 (empty response body)

### Authorization

[bearerAuthentication](../README.md#bearerAuthentication)

### HTTP request headers

- **Content-Type**: text/plain
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

