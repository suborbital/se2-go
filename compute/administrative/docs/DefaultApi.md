# \DefaultApi

All URIs are relative to *http://local.suborbital.network:8081/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetToken**](DefaultApi.md#GetToken) | **Get** /token/{environment}.{customerID}/{namespace}/{fnName} | 



## GetToken

> Token GetToken(ctx, environment, customerID, namespace, fnName).Execute()





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
    resp, r, err := api_client.DefaultApi.GetToken(context.Background(), environment, customerID, namespace, fnName).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetToken``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetToken`: Token
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetToken`: %v\n", resp)
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

Other parameters are passed through a pointer to a apiGetTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------





### Return type

[**Token**](Token.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

