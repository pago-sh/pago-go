# CustomerPortal.Downloadables

## Overview

### Available Operations

* [List](#list) - List Downloadables

## List

**Scopes**: `customer_portal:read` `customer_portal:write`

### Example Usage

<!-- UsageSnippet language="go" operationID="customer_portal:downloadables:list" method="get" path="/v1/customer-portal/downloadables/" -->
```go
package main

import(
	"context"
	pagogo "github.com/pago-sh/pago-go"
	"os"
	"github.com/pago-sh/pago-go/models/operations"
	"log"
)

func main() {
    ctx := context.Background()

    s := pagogo.New()

    res, err := s.CustomerPortal.Downloadables.List(ctx, operations.CustomerPortalDownloadablesListSecurity{
        CustomerSession: pagogo.Pointer(os.Getenv("PAGO_CUSTOMER_SESSION")),
    }, nil, pagogo.Pointer[int64](1), pagogo.Pointer[int64](10))
    if err != nil {
        log.Fatal(err)
    }
    if res.ListResourceDownloadableRead != nil {
        for {
            // handle items

            res, err = res.Next()

            if err != nil {
                // handle error
            }

            if res == nil {
                break
            }
        }
    }
}
```

### Parameters

| Parameter                                                                                                                                                   | Type                                                                                                                                                        | Required                                                                                                                                                    | Description                                                                                                                                                 |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                                                                       | [context.Context](https://pkg.go.dev/context#Context)                                                                                                       | :heavy_check_mark:                                                                                                                                          | The context to use for the request.                                                                                                                         |
| `security`                                                                                                                                                  | [operations.CustomerPortalDownloadablesListSecurity](../../models/operations/customerportaldownloadableslistsecurity.md)                                    | :heavy_check_mark:                                                                                                                                          | The security requirements to use for the request.                                                                                                           |
| `benefitID`                                                                                                                                                 | [*operations.CustomerPortalDownloadablesListQueryParamBenefitIDFilter](../../models/operations/customerportaldownloadableslistqueryparambenefitidfilter.md) | :heavy_minus_sign:                                                                                                                                          | Filter by benefit ID.                                                                                                                                       |
| `page`                                                                                                                                                      | `*int64`                                                                                                                                                    | :heavy_minus_sign:                                                                                                                                          | Page number, defaults to 1.                                                                                                                                 |
| `limit`                                                                                                                                                     | `*int64`                                                                                                                                                    | :heavy_minus_sign:                                                                                                                                          | Size of a page, defaults to 10. Maximum is 100.                                                                                                             |
| `opts`                                                                                                                                                      | [][operations.Option](../../models/operations/option.md)                                                                                                    | :heavy_minus_sign:                                                                                                                                          | The options for this request.                                                                                                                               |

### Response

**[*operations.CustomerPortalDownloadablesListResponse](../../models/operations/customerportaldownloadableslistresponse.md), error**

### Errors

| Error Type                    | Status Code                   | Content Type                  |
| ----------------------------- | ----------------------------- | ----------------------------- |
| apierrors.HTTPValidationError | 422                           | application/json              |
| apierrors.APIError            | 4XX, 5XX                      | \*/\*                         |